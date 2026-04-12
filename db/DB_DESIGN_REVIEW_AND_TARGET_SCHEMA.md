# WACAST Database Review and Target Schema

Dokumen ini mereview schema yang ada di `db/wecast.sql`, lalu menyesuaikannya dengan kebutuhan dashboard yang sudah dibuat agar hasil akhirnya:
- lebih ter-normalisasi
- konsisten naming-nya
- siap untuk scale
- lebih aman untuk query operasional dan analytics

## 1. Ringkasan Kondisi Schema Saat Ini

Tabel yang sudah ada di `wecast.sql`:
- `api_logs`
- `auto_response`
- `billing_plans`
- `broadcast_campaigns`
- `broadcast_messages`
- `broadcast_recipients`
- `contact`
- `devices`
- `groups`
- `lookup`
- `messages`
- `subscriptions`
- `system_settings`
- `users`
- `warming_pool`
- `warming_sessions`
- `webhooks`

## 2. Temuan Utama dari Schema Saat Ini

### A. Belum ada foreign key
Saat ini schema dump hanya punya primary key. Tidak ada:
- `FOREIGN KEY`
- `REFERENCES`
- `UNIQUE`
- index query penting

Dampak:
- integritas data lemah
- orphan records mudah terjadi
- join antar tabel rawan salah
- performa query dashboard akan turun saat data membesar

### B. Naming belum konsisten
Schema saat ini mencampur:
- camelCase: `userId`, `groupId`, `deviceId`, `nameBroadcast`
- snake_case: `created_at`, `last_seen`, `message_type`
- singular/plural tidak konsisten: `contact`, `groups`, `messages`

Rekomendasi:
- pakai `snake_case` konsisten
- pakai nama tabel plural konsisten
- pakai nama FK seragam seperti `user_id`, `device_id`, `group_id`

### C. Ada data yang bercampur antar concern
Contoh paling jelas:
- `users.otp_code` dan `users.otp_expired` seharusnya tidak ada di tabel `users`
- `users.idSubscribed` menduplikasi relasi yang seharusnya hidup di `subscriptions`
- `users.max_device` seharusnya turunan dari plan/subscription, bukan hardcoded di user kecuali memang override khusus
- `devices.session_data` terlalu berat bila dicampur langsung ke tabel device utama

### D. Normalisasi kontak dan grup belum benar
Saat ini tabel `contact` punya `groupId`.

Masalahnya:
- satu kontak hanya bisa ada di satu grup
- padahal dashboard group contact secara natural butuh many-to-many

Rekomendasi:
- `contacts`
- `contact_groups`
- `contact_group_members`

### E. Messaging belum cukup kaya untuk kebutuhan dashboard
Tabel `messages` saat ini terlalu minim untuk mendukung:
- message logs lengkap
- scheduled messages
- template usage
- message status progression
- broadcast lineage
- analytics/failure dashboard

Contoh gap:
- tidak ada `user_id`
- tidak ada `template_id`
- tidak ada `broadcast_id`
- tidak ada `scheduled_message_id`
- `direction` pakai `bit(1)` padahal sudah ada enum `enum_direction`
- belum ada `sent_at`, `delivered_at`, `failed_at`, `updated_at`

### F. Broadcast masih bisa dirapikan
Sekarang ada:
- `broadcast_campaigns`
- `broadcast_messages`
- `broadcast_recipients`

Ini sudah mendekati baik, tetapi:
- `broadcast_messages` terlihat seperti single payload per campaign
- jika memang satu campaign hanya satu message, payload bisa langsung disimpan di tabel campaign atau dipisah jadi `broadcast_contents` dengan alasan jelas

### G. Webhook masih terlalu sempit
Saat ini `webhooks` hanya punya:
- `deviceId`
- `webhookUrl`
- `secretKey`

Padahal dashboard butuh:
- scope per user/workspace
- daftar event yang diaktifkan
- delivery logs / retry logs

### H. Billing belum cukup untuk dashboard billing page
Yang ada baru:
- `billing_plans`
- `subscriptions`

Yang belum ada:
- `invoices`
- `invoice_items`

### H2. Settings dan onboarding belum sepenuhnya tercakup
Dashboard settings saat ini juga butuh data untuk:
- onboarding checklist / progress
- notification preferences
- bell notification history

Yang belum ada atau belum disebut eksplisit di schema existing:
- `onboarding_progress`
- `notification_settings`
- `notifications`
- `usage_quotas`

### I. Monitoring dashboard butuh tabel agregat
Dashboard punya halaman:
- usage statistics
- failure rate
- system health

Schema sekarang belum punya tabel agregat seperti:
- `daily_message_stats`
- `failure_records`
- `service_health_checks`
- `resource_usage_metrics`

### J. Tabel generik `lookup` sebaiknya dibatasi
`lookup` bisa berguna untuk enum dinamis, tetapi kalau dipakai terlalu luas akan jadi anti-pattern.

Rekomendasi:
- pakai enum PostgreSQL atau tabel spesifik bila domain-nya stabil
- jangan jadikan `lookup` tempat semua status bercampur

## 3. Rekomendasi Desain yang Lebih Normalized

## Domain 1: Identity and Auth

### `users`
Tabel utama user.

Kolom utama:
- `id` uuid pk
- `phone_number` varchar unique not null
- `full_name` varchar
- `email` varchar unique nullable
- `company_name` varchar nullable
- `timezone` varchar not null default `Asia/Jakarta`
- `is_verified` boolean not null default false
- `is_banned` boolean not null default false
- `is_api_enabled` boolean not null default false
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- `last_login_at` timestamptz nullable

### `otp_verifications`
Pisahkan OTP dari `users`.

Kolom utama:
- `id` uuid pk
- `user_id` uuid nullable fk `users.id`
- `phone_number` varchar not null
- `context` varchar not null
- `otp_code_hash` varchar not null
- `attempt_count` int not null default 0
- `expires_at` timestamptz not null
- `verified_at` timestamptz nullable
- `created_at` timestamptz not null

### `user_sessions`
Sesi login dashboard.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `session_token_hash` varchar not null
- `ip_address` inet nullable
- `user_agent` text nullable
- `expires_at` timestamptz not null
- `revoked_at` timestamptz nullable
- `created_at` timestamptz not null

## Domain 2: Billing and Subscription

### `billing_plans`
Tabel ini bisa dipertahankan, tetapi tambah metadata dasar.

Kolom utama:
- `id` uuid pk
- `name` varchar not null unique
- `price` numeric(12,2) not null
- `max_devices` int not null
- `max_messages_per_day` int nullable
- `features` jsonb not null default '{}'
- `is_active` boolean not null default true
- `created_at` timestamptz not null
- `updated_at` timestamptz not null

### `subscriptions`
Tabel ini tetap dipakai, tetapi jangan simpan relasi plan aktif juga di `users`.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `plan_id` uuid not null fk `billing_plans.id`
- `status` varchar not null
- `start_date` timestamptz not null
- `end_date` timestamptz nullable
- `renewal_date` timestamptz nullable
- `auto_renew` boolean not null default true
- `created_at` timestamptz not null
- `updated_at` timestamptz not null

### `usage_quotas`
Untuk dashboard kuota.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `subscription_id` uuid not null fk `subscriptions.id`
- `period_key` varchar not null
- `messages_used` int not null default 0
- `messages_limit` int nullable
- `devices_used` int not null default 0
- `devices_limit` int nullable
- `created_at` timestamptz not null
- `updated_at` timestamptz not null

### `invoices`
Untuk halaman billing.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `subscription_id` uuid nullable fk `subscriptions.id`
- `invoice_number` varchar not null unique
- `issue_date` date not null
- `due_date` date nullable
- `paid_at` timestamptz nullable
- `amount` numeric(12,2) not null
- `currency` varchar not null default 'IDR'
- `status` varchar not null
- `payment_method` varchar nullable
- `created_at` timestamptz not null

### `invoice_items`
Opsional untuk future-proof.

## Domain 3: Devices and WA Sessions

### `devices`
Refactor dari tabel `devices` sekarang.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `unique_name` varchar not null
- `display_name` varchar not null
- `phone_number` varchar nullable
- `status` varchar not null
- `platform` varchar nullable
- `wa_version` varchar nullable
- `battery_level` int nullable
- `last_seen_at` timestamptz nullable
- `connected_since` timestamptz nullable
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- unique `(user_id, unique_name)`

### `device_sessions`
Pindahkan session detail dari `devices.session_data` ke tabel terpisah.

Kolom utama:
- `id` uuid pk
- `device_id` uuid not null fk `devices.id`
- `session_blob` bytea nullable
- `session_status` varchar not null
- `started_at` timestamptz not null
- `ended_at` timestamptz nullable
- `last_restart_at` timestamptz nullable
- `created_at` timestamptz not null

### `device_qr_codes`
Untuk pairing QR.

### `device_metrics`
Untuk info/status/monitoring device.

## Domain 4: Contacts

### `contacts`
Refactor dari tabel `contact`.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `name` varchar not null
- `phone_number` varchar not null
- `label_id` uuid nullable fk `contact_labels.id`
- `additional_data` jsonb nullable
- `note` text nullable
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- `deleted_at` timestamptz nullable
- unique `(user_id, phone_number)`

### `contact_labels`
Baru. Dibutuhkan dashboard phonebook.

### `contact_groups`
Refactor dari tabel `groups`.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `name` varchar not null
- `description` text nullable
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- `deleted_at` timestamptz nullable

### `contact_group_members`
Tabel pivot yang wajib untuk normalisasi.

Kolom utama:
- `id` uuid pk
- `group_id` uuid not null fk `contact_groups.id`
- `contact_id` uuid not null fk `contacts.id`
- `created_at` timestamptz not null
- unique `(group_id, contact_id)`

### `blacklists`
Baru. Dibutuhkan langsung oleh dashboard blacklist.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `phone_number` varchar not null
- `reason` varchar nullable
- `blocked_at` timestamptz not null
- `unblocked_at` timestamptz nullable
- `created_at` timestamptz not null

## Domain 5: Messaging

### `messages`
Refactor besar dari tabel `messages` sekarang.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `device_id` uuid not null fk `devices.id`
- `direction` enum_direction not null
- `recipient_phone` varchar nullable
- `sender_phone` varchar nullable
- `message_type` varchar not null
- `content` text not null
- `status` varchar not null
- `error_log` text nullable
- `template_id` uuid nullable fk `message_templates.id`
- `broadcast_id` uuid nullable fk `broadcast_campaigns.id`
- `scheduled_message_id` uuid nullable fk `scheduled_messages.id`
- `created_at` timestamptz not null
- `sent_at` timestamptz nullable
- `delivered_at` timestamptz nullable
- `failed_at` timestamptz nullable
- `updated_at` timestamptz not null

Catatan penting:
- jangan pakai `bit(1)` untuk `direction`
- pakai enum `enum_direction` yang sudah ada
- `status_message` lebih baik diganti `status` varchar atau enum yang jelas

### `message_templates`
Baru. Dibutuhkan oleh dashboard templates.

### `scheduled_messages`
Baru. Dibutuhkan langsung oleh halaman scheduled.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `device_id` uuid not null fk `devices.id`
- `template_id` uuid nullable fk `message_templates.id`
- `group_id` uuid nullable fk `contact_groups.id`
- `recipient_mode` varchar not null
- `recipient_payload` jsonb nullable
- `message_content` text not null
- `scheduled_at` timestamptz not null
- `status` varchar not null
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- `executed_at` timestamptz nullable

Catatan normalisasi:
- `recipient_payload` JSONB boleh dipakai untuk tahap awal
- tetapi untuk desain yang benar-benar normalized dan queryable, sebaiknya tambah tabel turunan `scheduled_message_recipients`

### `scheduled_message_recipients`
Tabel child untuk recipient scheduled message agar tidak bergantung ke JSON.

Kolom utama:
- `id` uuid pk
- `scheduled_message_id` uuid not null fk `scheduled_messages.id`
- `contact_id` uuid nullable fk `contacts.id`
- `group_id` uuid nullable fk `contact_groups.id`
- `phone_number` varchar not null
- `status` varchar nullable
- `created_at` timestamptz not null
- unique `(scheduled_message_id, phone_number)`

### `broadcast_campaigns`
Tabel ini bisa dipertahankan dengan sedikit perapian.

Kolom utama yang disarankan:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `device_id` uuid not null fk `devices.id`
- `name` varchar not null
- `message_content` text nullable
- `template_id` uuid nullable fk `message_templates.id`
- `delay_seconds` int nullable
- `total_recipients` int not null default 0
- `processed_count` int not null default 0
- `success_count` int not null default 0
- `failed_count` int not null default 0
- `scheduled_at` timestamptz nullable
- `status` varchar not null
- `created_at` timestamptz not null
- `updated_at` timestamptz not null
- `deleted_at` timestamptz nullable

### `broadcast_messages`
Pilihan:
- jika campaign hanya satu pesan, tabel ini bisa dihapus dan isinya dipindah ke `broadcast_campaigns`
- jika nanti campaign mendukung multi-step payload, pertahankan tapi rename jadi `broadcast_contents`

### `broadcast_recipients`
Tabel ini tetap relevan.

Perbaikan kolom:
- `campaign_id` fk
- `contact_id` nullable fk
- `group_id` nullable fk
- `phone_number` varchar not null
- `status` varchar not null
- `sent_at` timestamptz nullable
- `error_message` text nullable
- `retry_count` int not null default 0
- index `(campaign_id, status)`

## Domain 6: Auto Response

### `auto_response_keywords`
Refactor dari `auto_response`.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `device_id` uuid nullable fk `devices.id`
- `keyword` varchar not null
- `response_text` text not null
- `is_active` boolean not null default true
- `created_at` timestamptz not null
- `updated_at` timestamptz not null

## Domain 7: Webhooks and API

### `webhooks`
Refactor dari tabel sekarang.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `device_id` uuid nullable fk `devices.id`
- `webhook_url` varchar not null
- `secret_key_hash` varchar not null
- `is_active` boolean not null default true
- `created_at` timestamptz not null
- `updated_at` timestamptz not null

### `webhook_event_subscriptions`
Baru. Dibutuhkan dashboard webhook events.

### `webhook_deliveries`
Baru. Untuk retry dan observability.

### `api_keys`
Baru. Dibutuhkan dashboard API keys.

Kolom utama:
- `id` uuid pk
- `user_id` uuid not null fk `users.id`
- `name` varchar not null
- `key_prefix` varchar not null
- `key_hash` varchar not null
- `last_used_at` timestamptz nullable
- `expires_at` timestamptz nullable
- `is_active` boolean not null default true
- `created_at` timestamptz not null
- `revoked_at` timestamptz nullable

### `api_logs`
Tabel ini bisa dipertahankan tetapi diperbaiki.

Perbaikan penting:
- `created_at` jangan `varchar`, harus `timestamptz`
- tambah FK ke `users` dan `devices`
- tambahkan index `created_at`

## Domain 8: Monitoring and Analytics

### `daily_message_stats`
Baru. Aggregated table untuk dashboard usage.

### `failure_records`
Baru. Untuk halaman failure rate.

### `service_health_checks`
Opsional. Untuk fitur system health bila mau diaktifkan lagi.

### `resource_usage_metrics`
Opsional. Untuk CPU/memory/disk.

## Domain 9: Notifications and Audit

### `notification_settings`
Baru. Dibutuhkan oleh page notifications.

### `notifications`
Baru. Dibutuhkan oleh bell dropdown/header.

### `audit_logs`
Baru. Untuk tracing perubahan penting.

## 4. Mapping Current Schema -> Target Schema

### Tabel yang bisa dipertahankan dengan refactor
- `billing_plans` -> tetap, rapikan kolom dan timestamp
- `subscriptions` -> tetap, tambah lifecycle columns dan FK
- `devices` -> tetap, pecah session blob dan tambah metadata
- `messages` -> tetap, tapi refactor besar
- `broadcast_campaigns` -> tetap, rapikan payload dan status
- `broadcast_recipients` -> tetap, rapikan FK dan status
- `api_logs` -> tetap, ganti timestamp dan tambahkan FK/index
- `webhooks` -> tetap, ubah scope jadi user-centric
- `warming_pool` -> tetap bila fitur warming memang core
- `warming_sessions` -> tetap bila fitur warming memang dipakai backend

### Tabel yang sebaiknya di-rename / di-normalisasi ulang
- `contact` -> `contacts`
- `groups` -> `contact_groups`
- `auto_response` -> `auto_response_keywords`

### Tabel yang sebaiknya ditambahkan
- `otp_verifications`
- `user_sessions`
- `usage_quotas`
- `invoices`
- `invoice_items`
- `device_sessions`
- `device_qr_codes`
- `device_metrics`
- `contact_labels`
- `contact_group_members`
- `blacklists`
- `message_templates`
- `scheduled_messages`
- `scheduled_message_recipients`
- `api_keys`
- `webhook_event_subscriptions`
- `webhook_deliveries`
- `daily_message_stats`
- `failure_records`
- `onboarding_progress`
- `notification_settings`
- `notifications`
- `audit_logs`

### Tabel/kolom yang sebaiknya dihapus atau dipindahkan
- `users.otp_code` -> pindah ke `otp_verifications`
- `users.otp_expired` -> pindah ke `otp_verifications`
- `users.idSubscribed` -> jangan simpan di user, pakai `subscriptions`
- `users.max_device` -> derived from plan/subscription kecuali override khusus
- `contact.groupId` -> ganti pivot `contact_group_members`
- `devices.session_data` -> pindah ke `device_sessions`
- `messages.direction bit(1)` -> ganti pakai enum

## 5. Index yang Wajib untuk Scale

Minimal index yang harus ada:

### Users
- unique index `users(phone_number)`
- unique index `users(email)` where email is not null

### Devices
- unique index `(user_id, unique_name)`
- index `(user_id, status)`
- index `(user_id, last_seen_at desc)`

### Contacts
- unique index `(user_id, phone_number)`
- index `(user_id, name)`

### Groups
- index `(user_id, name)`
- unique index `(group_id, contact_id)` pada pivot

### Messages
- index `(device_id, created_at desc)`
- index `(user_id, created_at desc)`
- index `(status, created_at desc)`
- index `(recipient_phone)`
- index `(broadcast_id)`
- index `(scheduled_message_id)`

### Broadcast
- index `(user_id, created_at desc)`
- index `(campaign_id, status)` pada recipients

### Webhooks / API
- index `(user_id, is_active)` pada `webhooks`
- index `(user_id, is_active)` pada `api_keys`
- index `(webhook_id, created_at desc)` pada delivery logs
- index `(created_at desc)` pada `api_logs`

### Analytics
- unique index `(user_id, device_id, stat_date)` pada `daily_message_stats`
- index `(user_id, occurred_at desc)` pada `failure_records`

## 6. Constraint dan Integritas Data yang Disarankan

- semua FK diberi `REFERENCES`
- gunakan `ON DELETE CASCADE` hanya untuk tabel turunan murni seperti pivot/log tertentu
- gunakan `ON DELETE RESTRICT` untuk resource penting seperti billing plan dan subscription history
- tambahkan check constraint pada `battery_level between 0 and 100`
- tambahkan check constraint untuk angka non-negatif seperti retry, quota, count
- simpan nomor HP dalam format yang konsisten
- secret/API key disimpan hashed, bukan plaintext

## 7. Status Normalisasi

### Sudah lumayan dekat
- billing plan vs subscription sudah terpisah
- broadcast recipient sudah dipisah dari campaign
- webhook sudah dipisah dari device

### Belum normal / harus diperbaiki
- OTP masih ditempel ke `users`
- kontak masih one-group-only
- user menyimpan data subscription turunan
- messages terlalu padat tapi informasinya justru kurang lengkap
- scheduled message recipient belum dipisah jika ingin normalized penuh
- event/config/log webhook belum dipisah
- tabel dashboard settings/notifications belum ada
- onboarding progress belum punya persistence table

## 8. Minimal Target Schema yang Saya Sarankan

Kalau targetnya pragmatic tapi tetap scalable, saya sarankan minimal final schema berisi:

1. `users`
2. `otp_verifications`
3. `user_sessions`
4. `billing_plans`
5. `subscriptions`
6. `usage_quotas`
7. `invoices`
8. `devices`
9. `device_sessions`
10. `device_qr_codes`
11. `contacts`
12. `contact_labels`
13. `contact_groups`
14. `contact_group_members`
15. `blacklists`
16. `messages`
17. `message_templates`
18. `scheduled_messages`
19. `scheduled_message_recipients`
20. `broadcast_campaigns`
21. `broadcast_recipients`
22. `auto_response_keywords`
23. `api_keys`
24. `webhooks`
25. `webhook_event_subscriptions`
26. `webhook_deliveries`
27. `api_logs`
28. `daily_message_stats`
29. `failure_records`
30. `onboarding_progress`
31. `notification_settings`
32. `notifications`
33. `audit_logs`

## 9. Kesimpulan

Schema `wecast.sql` saat ini masih cukup bagus sebagai starting point backend operasional, tetapi belum cukup rapi untuk langsung menopang dashboard modern secara penuh.

Masalah paling besar ada di:
- tidak adanya foreign key dan index
- normalisasi contacts/groups yang belum benar
- auth/OTP masih menempel di user
- progress onboarding belum punya persistence table
- scheduled message recipient belum dipisah jika ingin normalized penuh
- billing, notifications, API keys, template, scheduled message, dan analytics aggregate belum lengkap
- beberapa tipe data dan nama kolom masih belum konsisten

Arah terbaik bukan membuang schema lama, tetapi:
- pertahankan tabel core yang sudah benar arahnya
- refactor naming dan relation
- pecah concern yang masih bercampur
- tambahkan tabel khusus untuk dashboard-oriented features

Dengan struktur target di atas, database akan jauh lebih aman, lebih mudah di-query, dan lebih siap scale untuk multi-user, multi-device, dan analytics dashboard.

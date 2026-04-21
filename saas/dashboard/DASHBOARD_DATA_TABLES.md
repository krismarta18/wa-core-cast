# Dashboard Data Tables

Dokumen ini merangkum kebutuhan tabel database yang dibutuhkan untuk menampung data seluruh page, komponen, dan fitur di web dashboard pada folder `src/app` dan komponen pendukung layout/dashboard.

## Cakupan Fitur

### Auth
- Login via nomor HP
- Register akun
- OTP verification

### Dashboard Home
- Ringkasan paket aktif
- Kuota pesan
- Slot device
- Status device aktif
- Pesan terbaru
- Statistik harian

### Device Management
- Status device
- Pairing QR
- Device info
- Session management

### Messaging
- New message
- Broadcast
- Scheduled messages
- Message logs

### Contact Management
- Phonebook
- Contact groups
- Blacklist

### Auto Response
- Keyword auto-reply
- Message templates

### API & Integration
- API keys
- Webhooks

### Monitoring
- Usage statistics
- Failure rate
- System health

Catatan:
`System health` masih ada di source code, walau sudah di-take out dari navigasi.

### Settings
- Onboarding
- Profile
- Billing
- Notifications

## Kebutuhan Tabel Utama

## 1. Identity dan Auth

### `users`
Menyimpan data akun utama dashboard.

Kolom utama:
- `id`
- `full_name`
- `phone_number`
- `email`
- `company_name`
- `timezone`
- `password_hash` atau nullable bila full OTP auth
- `is_active`
- `last_login_at`
- `created_at`
- `updated_at`

### `otp_verifications`
Menyimpan OTP untuk login/register/verify action.

Kolom utama:
- `id`
- `user_id` nullable
- `phone_number`
- `context` (`login`, `register`, `reset`, dll)
- `otp_code`
- `attempt_count`
- `expires_at`
- `verified_at`
- `created_at`

### `user_sessions`
Menyimpan sesi login dashboard.

Kolom utama:
- `id`
- `user_id`
- `access_token` atau `session_token`
- `refresh_token` nullable
- `ip_address`
- `user_agent`
- `expired_at`
- `created_at`
- `revoked_at`

## 2. Subscription dan Billing

### `plans`
Master plan subscription.

Kolom utama:
- `id`
- `name`
- `price_monthly`
- `message_quota`
- `device_limit`
- `feature_flags` JSON
- `is_active`
- `created_at`
- `updated_at`

### `subscriptions`
Plan aktif per user/workspace.

Kolom utama:
- `id`
- `user_id`
- `plan_id`
- `status`
- `start_date`
- `end_date`
- `renewal_date`
- `auto_renew`
- `created_at`
- `updated_at`

### `usage_quotas`
Tracking kuota bulanan user.

Kolom utama:
- `id`
- `user_id`
- `subscription_id`
- `period_month`
- `messages_used`
- `messages_limit`
- `devices_used`
- `devices_limit`
- `created_at`
- `updated_at`

### `invoices`
Invoice billing seperti yang tampil di halaman Billing.

Kolom utama:
- `id`
- `user_id`
- `subscription_id`
- `invoice_number`
- `issue_date`
- `due_date`
- `paid_at` nullable
- `amount`
- `currency`
- `status`
- `payment_method` nullable
- `created_at`

### `invoice_items`
Rincian item per invoice bila nanti invoice lebih kompleks.

Kolom utama:
- `id`
- `invoice_id`
- `description`
- `qty`
- `unit_price`
- `total_price`

## 3. Device dan Session WhatsApp

### `devices`
Entitas utama seluruh device WhatsApp user.

Kolom utama:
- `id`
- `user_id`
- `name`
- `phone_number`
- `status` (`connected`, `disconnected`, `pending_qr`, dll)
- `platform`
- `wa_version`
- `battery_level`
- `connected_since`
- `last_seen_at`
- `last_disconnect_reason` nullable
- `created_at`
- `updated_at`

### `device_sessions`
Menyimpan state session device untuk login/logout/restart.

Kolom utama:
- `id`
- `device_id`
- `session_identifier`
- `status`
- `started_at`
- `ended_at` nullable
- `restart_count`
- `last_restart_at` nullable
- `logout_at` nullable

### `device_qr_codes`
Menyimpan QR pairing.

Kolom utama:
- `id`
- `device_id` nullable
- `user_id`
- `qr_string`
- `qr_image_url` atau `qr_image_base64`
- `status` (`pending`, `scanned`, `expired`)
- `generated_at`
- `expired_at`

### `device_metrics`
Metrik operasional per device, untuk status/info/monitoring.

Kolom utama:
- `id`
- `device_id`
- `uptime_seconds`
- `messages_sent_count`
- `messages_received_count`
- `success_rate`
- `recorded_at`

## 4. Contact Management

### `contact_labels`
Master label phonebook.

Kolom utama:
- `id`
- `user_id`
- `name`
- `color` nullable
- `created_at`

### `contacts`
Phonebook utama.

Kolom utama:
- `id`
- `user_id`
- `name`
- `phone_number`
- `label_id` nullable
- `note` nullable
- `created_at`
- `updated_at`

### `contact_groups`
Group kontak untuk broadcast.

Kolom utama:
- `id`
- `user_id`
- `name`
- `description` nullable
- `created_at`
- `updated_at`

### `contact_group_members`
Pivot members dari group kontak.

Kolom utama:
- `id`
- `group_id`
- `contact_id`
- `created_at`

### `blacklists`
Nomor yang diblokir.

Kolom utama:
- `id`
- `user_id`
- `phone_number`
- `reason`
- `blocked_at`
- `unblocked_at` nullable
- `created_at`

## 5. Messaging

### `messages`
Tabel utama seluruh pesan outgoing/incoming.

Kolom utama:
- `id`
- `user_id`
- `device_id`
- `direction` (`inbound`, `outbound`)
- `message_type` (`single`, `broadcast`, `scheduled`, `auto_reply`, `otp`)
- `recipient_phone` nullable
- `sender_phone` nullable
- `content`
- `status` (`pending`, `sent`, `delivered`, `read`, `failed`, `cancelled`)
- `scheduled_at` nullable
- `sent_at` nullable
- `failed_at` nullable
- `failure_reason` nullable
- `created_at`
- `updated_at`

### `message_logs`
Opsional jika ingin pisahkan log operasional dari message entity utama.

Kolom utama:
- `id`
- `message_id`
- `device_id`
- `status`
- `status_detail` nullable
- `logged_at`

Catatan:
Jika `messages` sudah cukup detail, tabel ini bisa tidak diperlukan.

### `broadcasts`
Header / job untuk broadcast.

Kolom utama:
- `id`
- `user_id`
- `device_id`
- `title` nullable
- `message_content`
- `template_id` nullable
- `delay_seconds`
- `status` (`draft`, `queued`, `sending`, `completed`, `failed`, `cancelled`)
- `total_recipients`
- `sent_count`
- `failed_count`
- `created_at`
- `started_at` nullable
- `completed_at` nullable

### `broadcast_recipients`
Detail recipient tiap broadcast.

Kolom utama:
- `id`
- `broadcast_id`
- `contact_id` nullable
- `group_id` nullable
- `phone_number`
- `status`
- `sent_at` nullable
- `failed_at` nullable
- `failure_reason` nullable

### `scheduled_messages`
Jadwal pengiriman pesan.

Kolom utama:
- `id`
- `user_id`
- `device_id`
- `template_id` nullable
- `group_id` nullable
- `recipient_mode` (`single`, `group`, `bulk`)
- `recipient_payload` JSON
- `message_content`
- `scheduled_at`
- `status`
- `created_at`
- `updated_at`
- `executed_at` nullable

## 6. Auto Response dan Template

### `message_templates`
Template pesan reusable.

Kolom utama:
- `id`
- `user_id`
- `name`
- `category`
- `content`
- `used_count`
- `created_at`
- `updated_at`

### `auto_response_keywords`
Aturan keyword auto-reply.

Kolom utama:
- `id`
- `user_id`
- `keyword`
- `response_text`
- `is_active`
- `created_at`
- `updated_at`

### `auto_response_logs`
Log trigger keyword auto-reply.

Kolom utama:
- `id`
- `keyword_id`
- `message_id` nullable
- `triggered_by_phone`
- `matched_keyword`
- `response_sent`
- `created_at`

## 7. API Keys dan Webhooks

### `api_keys`
API key milik user.

Kolom utama:
- `id`
- `user_id`
- `name`
- `key_prefix`
- `key_hash`
- `last_used_at` nullable
- `expires_at` nullable
- `is_active`
- `created_at`
- `revoked_at` nullable

### `webhooks`
Konfigurasi endpoint webhook.

Kolom utama:
- `id`
- `user_id`
- `endpoint_url`
- `secret`
- `is_active`
- `created_at`
- `updated_at`

### `webhook_event_subscriptions`
Event apa saja yang diaktifkan per webhook.

Kolom utama:
- `id`
- `webhook_id`
- `event_key`
- `is_enabled`

### `webhook_deliveries`
Log pengiriman webhook.

Kolom utama:
- `id`
- `webhook_id`
- `event_key`
- `payload` JSON
- `attempt`
- `http_status` nullable
- `response_body` nullable
- `status`
- `sent_at` nullable
- `created_at`

## 8. Monitoring dan Analytics

### `daily_message_stats`
Statistik agregat harian untuk halaman usage.

Kolom utama:
- `id`
- `user_id`
- `device_id` nullable
- `stat_date`
- `sent_count`
- `failed_count`
- `delivered_count` nullable
- `success_rate` nullable
- `created_at`

### `failure_records`
Data kegagalan pengiriman untuk halaman failure rate.

Kolom utama:
- `id`
- `user_id`
- `device_id`
- `message_id` nullable
- `recipient_phone`
- `failure_type`
- `failure_reason`
- `occurred_at`
- `created_at`

### `service_health_checks`
Untuk halaman system health bila fitur ini ingin diaktifkan lagi.

Kolom utama:
- `id`
- `service_name`
- `status`
- `latency_ms`
- `uptime_percent`
- `last_incident_at` nullable
- `checked_at`

### `resource_usage_metrics`
CPU, memory, disk metrics.

Kolom utama:
- `id`
- `metric_name`
- `metric_value`
- `metric_unit`
- `recorded_at`

## 9. Notifications dan In-App Alerts

### `notification_settings`
Preferensi notifikasi user.

Kolom utama:
- `id`
- `user_id`
- `event_key`
- `email_enabled`
- `in_app_enabled`
- `created_at`
- `updated_at`

### `notifications`
Data notif in-app yang tampil di header/bell.

Kolom utama:
- `id`
- `user_id`
- `type` (`warning`, `alert`, `success`, `info`)
- `title`
- `body`
- `is_read`
- `read_at` nullable
- `created_at`

## 10. Audit dan Activity

### `audit_logs`
Audit trail perubahan penting.

Kolom utama:
- `id`
- `user_id`
- `action_type`
- `resource_type`
- `resource_id`
- `metadata` JSON nullable
- `ip_address` nullable
- `user_agent` nullable
- `created_at`

## Tabel Minimum yang Paling Penting

Kalau mau mulai dari struktur minimum dulu, tabel inti yang paling wajib adalah:

1. `users`
2. `otp_verifications`
3. `user_sessions`
4. `plans`
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
17. `broadcasts`
18. `broadcast_recipients`
19. `scheduled_messages`
20. `message_templates`
21. `auto_response_keywords`
22. `api_keys`
23. `webhooks`
24. `webhook_event_subscriptions`
25. `webhook_deliveries`
26. `daily_message_stats`
27. `failure_records`
28. `notification_settings`
29. `notifications`
30. `audit_logs`

## Tabel Optional / Tahap Lanjut

Bisa ditambahkan belakangan jika memang diperlukan:

- `invoice_items`
- `device_metrics`
- `message_logs`
- `auto_response_logs`
- `service_health_checks`
- `resource_usage_metrics`

## Relasi Penting

Relasi utama yang perlu dijaga:

- 1 user punya banyak device
- 1 user punya banyak contact
- 1 group punya banyak contact melalui `contact_group_members`
- 1 user punya banyak message, broadcast, scheduled message
- 1 device bisa dipakai oleh banyak message
- 1 template bisa dipakai oleh banyak message/broadcast/scheduled message
- 1 user punya banyak API key dan webhook
- 1 webhook punya banyak event subscription dan delivery log
- 1 user punya banyak notification dan audit log

## Saran Implementasi

- Gunakan UUID untuk primary key pada tabel utama
- Simpan nomor telepon dalam format E.164
- Simpan secret / API key dalam bentuk hash, bukan plaintext
- Gunakan enum atau lookup table untuk status yang stabil
- Tabel analytics sebaiknya agregat harian agar dashboard cepat
- Pisahkan tabel operasional (`messages`) dan tabel agregat (`daily_message_stats`)
- `System health` bisa dianggap optional karena saat ini sudah dihapus dari navigasi dashboard

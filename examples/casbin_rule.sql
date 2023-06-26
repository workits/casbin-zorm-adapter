CREATE TABLE `casbin_rule` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `ptype` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v0` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v1` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v2` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v3` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v4` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `v5` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
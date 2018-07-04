Create TABLE bitfinex_order (
`id` int(11) NOT NULL AUTO_INCREMENT,
`order_id` VARCHAR(50) NOT NULL ,
`unique_id` VARCHAR(50) NOT NULL UNIQUE ,
`coin` VARCHAR(20) NOT NULL,
`market_coin`VARCHAR(20) NOT NULL,
`price` double NOT NULL,
`qty` double NOT NULL,
`action` VARCHAR(20) NOT NULL,
`expect_profit` double NOT NULL,
`actual_profit` double NOT NULL DEFAULT -1,
`expect_profit_rate` double NOT NULL,
`actual_price` double NOT NULL DEFAULT -1,
`expect_price` double NOT NULL,
`skid` double NOT NULL,
`status` VARCHAR(10) NOT NULL DEFAULT "",
`last_modify` bigint NOT NULL,
`created_ts` bigint NOT NULL ,
PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


Create TABLE binance_order (
`id` int(11) NOT NULL AUTO_INCREMENT,
`order_id` VARCHAR(50) NOT NULL ,
`unique_id` VARCHAR(50) NOT NULL UNIQUE ,
`coin` VARCHAR(20) NOT NULL,
`market_coin`VARCHAR(20) NOT NULL,
`price` double NOT NULL,
`qty` double NOT NULL,
`action` VARCHAR(20) NOT NULL,
`expect_profit` double NOT NULL,
`actual_profit` double NOT NULL DEFAULT -1,
`expect_profit_rate` double NOT NULL,
`actual_price` double NOT NULL DEFAULT -1,
`expect_price` double NOT NULL,
`skid` double NOT NULL,
`status` VARCHAR(10) NOT NULL DEFAULT "",
`last_modify` bigint NOT NULL,
`created_ts` bigint NOT NULL ,
PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

Create TABLE huobi_order (
`id` int(11) NOT NULL AUTO_INCREMENT,
`order_id` VARCHAR(50) NOT NULL ,
`unique_id` VARCHAR(50) NOT NULL UNIQUE ,
`coin` VARCHAR(20) NOT NULL,
`market_coin`VARCHAR(20) NOT NULL,
`price` double NOT NULL,
`qty` double NOT NULL,
`action` VARCHAR(20) NOT NULL,
`expect_profit` double NOT NULL,
`actual_profit` double NOT NULL DEFAULT -1,
`expect_profit_rate` double NOT NULL,
`actual_price` double NOT NULL DEFAULT -1,
`expect_price` double NOT NULL,
`skid` double NOT NULL,
`status` VARCHAR(10) NOT NULL DEFAULT "",
`last_modify` bigint NOT NULL,
`created_ts` bigint NOT NULL ,
PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
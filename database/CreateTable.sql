USE fizzbuzz;
CREATE TABLE request_history (
                           `id` INT AUTO_INCREMENT primary key NOT NULL,
                           `int1` int(11) NOT NULL,
                           `int2` int(11) NOT NULL,
                           `limit` int(11) NOT NULL,
                           `str1`  varchar(250) NOT NULL,
                           `str2` varchar(250) NOT NULL,
                           `count` int(11) NOT NULL DEFAULT 1
);
CREATE DATABASE fizzbuzz_test;
USE fizzbuzz_test;
CREATE TABLE request_history (
                                 `id` INT AUTO_INCREMENT primary key NOT NULL,
                                 `int1` int(11) NOT NULL,
                                 `int2` int(11) NOT NULL,
                                 `limit` int(11) NOT NULL,
                                 `str1`  varchar(250) NOT NULL,
                                 `str2` varchar(250) NOT NULL,
                                 `count` int(11) NOT NULL DEFAULT 1
);

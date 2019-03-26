CREATE TABLE orders
(
	id int NOT NULL AUTO_INCREMENT,
	exchangeAddress varchar(42) NOT NULL, 
	makerAddress varchar(42) NOT NULL, 
	takerAddress varchar(42) NOT NULL, 
	feeRecipientAddress char(42) NOT NULL,
	senderAddress varchar(42) NOT NULL,
	makerAssetAmount varchar(255) NOT NULL,
	takerAssetAmount varchar(255) NOT NULL,
	makerAssetFilledAmount varchar(255) NOT NULL,
	takerAssetFilledAmount varchar(255) NOT NULL,
	makerFee varchar(255) NOT NULL,
	takerFee varchar(255) NOT NULL,
	expirationTimeSeconds varchar(10) NOT NULL,
	salt varchar(255) NOT NULL,
	makerAssetData varchar(255) NOT NULL,
	takerAssetData varchar(255) NOT NULL,
	signature varchar(255) NOT NULL,
	PRIMARY KEY (id)
) default charset = utf8, ENGINE=InnoDB;
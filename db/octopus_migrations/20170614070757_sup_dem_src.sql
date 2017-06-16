
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE exchange.sup_dem_src CHANGE imp_in_count ad_in_count INT(11);
ALTER TABLE exchange.sup_dem_src CHANGE win_count ad_out_count INT(11);
ALTER TABLE exchange.sup_dem_src CHANGE win_bid ad_out_bid INT(11) ;

ALTER TABLE exchange.sup_src CHANGE imp_out_count ad_out_count INT(11) ;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE exchange.sup_dem_src DROP ad_in_count;
ALTER TABLE exchange.sup_dem_src DROP ad_out_count;
ALTER TABLE exchange.sup_dem_src DROP ad_out_bid;

ALTER TABLE exchange.sup_src DROP ad_in_count;







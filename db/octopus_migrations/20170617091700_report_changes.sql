
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE exchange.demand_report CHANGE imp_in_count ad_in_count INT(11);
ALTER TABLE exchange.demand_report CHANGE win_count ad_out_count INT(11);
ALTER TABLE exchange.demand_report CHANGE win_bid ad_out_bid INT(11) ;
ALTER TABLE exchange.supplier_report CHANGE impression_out ad_out_count INT(11) ;
ALTER TABLE exchange.supplier_report CHANGE impression_in impression_in_count INT(11) ;
ALTER TABLE exchange.supplier_report CHANGE delivered_impression delivered_count INT(11) ;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE exchange.demand_report DROP ad_in_count;
ALTER TABLE exchange.demand_report DROP ad_out_count;
ALTER TABLE exchange.demand_report DROP ad_out_bid;
ALTER TABLE exchange.supplier_report DROP ad_out_count;
ALTER TABLE exchange.supplier_report DROP delivered_count;
ALTER TABLE exchange.supplier_report DROP impression_in_count;












```sql
select reveal_txid, file_name from website.tb_inscribe_data where order_id in (select order_id from website.tb_inscribe_order where order_status != 'PAYTIMEOUT' and order_status != 'MINTING'  and `version`=137 );

mysql -uroot -pFUnxy7jdfYsxkdfs -e "select reveal_txid, file_name from website.tb_inscribe_data where order_id in (select order_id from website.tb_inscribe_order where order_status != 'PAYTIMEOUT' and order_status != 'MINTING'  and version=137 );
" > fish.txt

select receive_address, count  from tb_inscribe_order where order_status != 'PAYTIMEOUT' and `version`=137;


select receive_address, count, order_status from tb_inscribe_order where order_status != 'PAYTIMEOUT' and order_status != 'MINTING' and `version`=137;

select sum(count) from tb_inscribe_order where order_status != 'PAYTIMEOUT' and order_status != 'MINTING' and `version`=137;



select sum(count) from tb_inscribe_order where order_status != 'PAYTIMEOUT' and `version`=137;

select sum(mint_limit) from tb_waitlist;
```
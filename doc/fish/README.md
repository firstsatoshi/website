



```sql
select reveal_txid from tb_inscribe_data where order_id in (select order_id from tb_inscribe_order where order_status != 'PAYTIMEOUT' and receive_address='bc1psdkup70gwzjz9mljt9jtkm52z00gvzs2nwxvtapcxq9065mwdqkqpgsuyx');



select receive_address, count from tb_inscribe_order where order_status != 'PAYTIMEOUT' and `version`=137;



select sum(count) from tb_inscribe_order where order_status != 'PAYTIMEOUT' and `version`=137;

select sum(mint_limit) from tb_waitlist;
```
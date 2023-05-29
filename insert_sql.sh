
# docker cp doc/website-v1.0.sql  mysql:/root/
# docker exec -it mysql mysql -u root -pFUnxy7jdfYsxkdfs -D website < doc/website-v1.0.sql

docker exec -it mysql mysql -u root -pFUnxy7jdfYsxkdfs -e "INSERT INTO website.tb_blindbox_event (event_name,event_description,btc_price,is_active,payment_coin,supply,avail,only_whitelist,start_time,end_time,create_time,update_time) VALUES('Bitcoin Eagle','This is Bitcoin Eagle NFT mint',123456,1,'BTC',1000,1000,0,'2023-05-27 16:28:39','2024-05-27 16:28:39','2023-05-27 16:28:39','2023-05-27 16:28:48');"



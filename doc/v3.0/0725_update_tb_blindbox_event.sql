ALTER table tb_blindbox_event ADD COLUMN `custome_mint` tinyint(1) DEFAULT '0';



INSERT INTO website.tb_blindbox_event
(id, event_name, event_endpoint, event_description, detail, avatar_img_url, background_img_url, roadmap_description, roadmap_list, website_url, whitepaper_url, twitter_url, discord_url, price_sats, is_active, is_display, payment_token, current_mint_plan_index, mint_plan_list, img_url_list, average_image_bytes, supply, avail, lock_count, mint_limit, only_whitelist, start_time, end_time, create_time, update_time, custome_mint)
VALUES(2, 'BitcoinFish', 'bitcoinfish', 'BitcoinFish', 'BitcoinFish is from China''s earliest PFP (Profile Picture Project) project, CryptoFish, has now set sail onto the most venerable Bitcoin blockchain, adopting the latest Ordinals recursive inscription to conduct an unprecedented custom engraving. <br><br> Even fish have dreams, and this group of non-laid-back fish dreams of becoming pioneers in the digital art industry, leading the development of China''s digital art sector.<br><br>CryptoFish is China''s first native digital encrypted IP, consisting of 10,000 pixel-style fish avatars and 101 special fish. Each fish avatar is a unique digital collectible permanently stored on IPFS and minted according to the ERC721 contract standard on the blockchain. As China''s oldest PFP project, Crypto Fish has brought together a group of art enthusiasts, NFT collectors, and holds widespread influence in the Chinese NFT community, making it the Cryptopunks of China.<br><br>Now, BitcoinFish thanks to the application by the CryptoFish DAO Foundation, all 95 elements of the fish have been minted using Ordinals onto Bitcoin blockchain. Innovatively, we have introduced the latest recursive inscription, allowing users to customize their own engraved BitcoinFishes by choosing different elements. Each fish is personally selected and engraved by the user, making it unique, permanently stored on the Bitcoin, and ensuring that no previously combined fish elements are reused, making each BitcoinFish completely new and one-of-a-kind.<br><br>The holders of CryptoFish hope that it can go beyond borders and enter broader overseas and native public chain markets, allowing crypto collectors and enthusiasts worldwide to enjoy and purchase original PFP projects with elements from China – the BitcoinFish. This will also leave a profound mark of the Fish''s belief and cultural attributes, symbolizing the never laid-back spirit, on the Bitcoin blockchain, lasting for eternity!', 'https://static.nft.cc/bitfash//bitfish1.png', 'https://static.nft.cc/bitfash//bitfish1.png', 'roadmap', 'whitelist mint; public mint; maketplace launching  ', 'https://www.nft.cc/collection/bitcoinfish', 'https://www.nft.cc/collection/bitcoinfish', 'https://www.nft.cc/collection/bitcoinfish', 'https://www.nft.cc', 10000, 1, 1, 'BTC', 0, 'Whitelist Mint, 2023-07-31', 'https://static.nft.cc/bitfash/bitfish1.png;https://static.nft.cc/bitfash/bitfish2.png;https://static.nft.cc/bitfash/bitfish3.png;https://static.nft.cc/bitfash/bitfish4.png;https://static.nft.cc/bitfash/bitfish5.png', 500, 1000, 1000, 0, 1, 1, '2023-07-25 17:49:32', '2023-07-25 17:49:32', '2023-07-25 17:49:32', '2023-07-28 11:37:20', 1);

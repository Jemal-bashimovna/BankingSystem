SELECT EXISTS(SELECT 1 FROM accounts WHERE id=1);
INSERT INTO transactions (account_id, amount, transaction_type) VALUES (2, 200.00, 'account-withdraw');

UPDATE accounts SET balance=500.00 WHERE id=2;

 update accounts set balance=55.00 where id=1;
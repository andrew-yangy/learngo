create schema if not exists public;

create type order_status as enum (
    'created', 'accepted', 'pending',
    'cancelled', 'partially_filled', 'filled',
    'expired', 'rejected', 'confirmed'
);

create type status_enum as enum (
    'pending', 'accepted', 'rejected', 'confirmed', 'rolled_back', 'withdrawn'
);

create type event_type_enum as enum (
    'deposit', 'mint', 'withdrawal', 'transfer',
    'full_withdrawal', 'false_full_withdrawal', 'settlement', 'order'
);

create type event_status_enum as enum
(
    'pending', 'accepted', 'rejected', 'confirmed',
    'rolled_back', 'withdrawn', 'created', 'cancelled',
    'partially_filled', 'filled', 'expired'
);

create type transaction_enum as enum (
    'deposit', 'mint', 'withdrawal', 'transfer', 'full_withdrawal', 'false_full_withdrawal'
);

create type token_enum as enum (
    'ETH', 'ERC20', 'ERC721'
);

create table if not exists accounts (
    stark_key varchar primary key,
    ether_key varchar,
    nonce bigint check (nonce >= 0),
    tx_hash varchar(66),
    unique ( stark_key, ether_key )
);

create table if not exists contracts (
    contract_address varchar(42),
    owner varchar,
    asset_info varchar(74) not null, -- ETH: 4 bytes, ERC20/721: 36 bytes
    asset_type varchar primary key, -- uint
    unique (contract_address, owner, asset_info, asset_type)
);

create table if not exists tokens (
    id varchar primary key,
    type token_enum,
    contract_address varchar(42) not null,
    ticker_symbol varchar(5) not null,
    quantum bigint check (quantum >= 0),
    token_id varchar,
    client_token_id varchar,
    blueprint varchar
);

create table if not exists vaults (
    id integer primary key check (id >= 0 and id < 2147483648), -- 31 bit unsigned max vault id
    stark_key varchar references accounts (stark_key),
    asset_id varchar references tokens (id),
    quantized_balance bigint check (quantized_balance >= 0),
    unique (stark_key, asset_id)
);

create table if not exists orders (
    id bigint primary key,
    stark_key varchar references accounts (stark_key),
    status order_status,
    buy_vault_id integer references vaults (id),
    buy_quantity numeric(64, 0) check (buy_quantity > 0),
    sell_vault_id integer references vaults (id),
    sell_quantity numeric(64, 0) check (sell_quantity > 0),
    nonce bigint check (nonce >= 0),
    stark_signature varchar,
    creation_time timestamp default current_timestamp,
    expiration_time timestamp,
    unique (stark_key, nonce)
);

create table if not exists sequences (
    id bigint primary key check (id >= 0)
);

create table if not exists transactions (
    id bigint primary key references sequences (id),
    batch_id bigint check (batch_id >= -1),
    transaction_type transaction_enum,
    status status_enum,
    from_vault_id integer references vaults (id),
    to_vault_id integer references vaults (id),
    quantity numeric(64, 0) check (quantity >= 0),
    nonce bigint check (nonce >= 0),
    stark_signature varchar,
    tx_hash varchar(66),
    creation_time timestamp default current_timestamp,
    expiration_time timestamp
);

create table if not exists settlements (
    id bigint primary key references sequences (id),
    batch_id bigint check (batch_id >= -1),
    status status_enum,

    party_a_order_id bigint references orders (id),
    party_a_bought_vault_id integer references vaults (id),
    party_a_bought_quantity numeric(64, 0) check (party_a_bought_quantity >= 0),
    party_a_sold_vault_id integer references vaults (id),
    party_a_sold_quantity numeric(64, 0) check (party_a_sold_quantity >= 0),

    party_b_order_id bigint references orders (id),
    party_b_bought_vault_id integer references vaults (id),
    party_b_bought_quantity numeric(64, 0) check (party_b_bought_quantity >= 0),
    party_b_sold_vault_id integer references vaults (id),
    party_b_sold_quantity numeric(64, 0) check (party_b_sold_quantity >= 0),

    tx_hash varchar(66),
    creation_time timestamp default current_timestamp
);

create table if not exists events (
    id serial check (id > 0),
    type event_type_enum,
    status event_status_enum,
    reference_id bigint,
    batch_id bigint check (batch_id >= -1),
    creation_time timestamp default current_timestamp,
    reason varchar,
    context varchar,
    primary key (id)
);

create table if not exists transitions (
    id bigint primary key,
    batch_id bigint check (batch_id >= -1) DEFAULT -1,
    type event_type_enum,
    status event_status_enum,
    body jsonb default '{}'::jsonb,
    submitted_at timestamp default current_timestamp,
    successful boolean default false
);

-- Insert an asset to represent Ethereum in our system
INSERT INTO tokens (id, type, quantum, ticker_symbol, contract_address, token_id, client_token_id, blueprint) VALUES (
    '0x02705737cd248ac819034b5de474c8f0368224f72a0fda9e031499d519992d9e',
    'ETH',
    100000000,
    'ETH',
    '',
    '',
    '',
    ''
);

CREATE INDEX accounts_ether_key_index ON accounts(ether_key);

ALTER TABLE transactions DROP constraint transactions_from_vault_id_fkey;
ALTER TABLE transactions DROP constraint transactions_to_vault_id_fkey;

ALTER TABLE events ADD CONSTRAINT events_unique UNIQUE (type, status, reference_id, batch_id, reason);

ALTER TABLE transactions ALTER COLUMN batch_id SET DEFAULT -1;
ALTER TABLE settlements  ALTER COLUMN batch_id SET DEFAULT -1;

create table if not exists atrx (
    id bigint primary key,
    reason_code varchar not null,
    reason_msg varchar not null,
    tx jsonb default '{}'::jsonb,
    alt_txs jsonb default '[]'::jsonb
);


ALTER TABLE tokens ADD CONSTRAINT token_field_disjunction CHECK
((type = 'ETH' AND (contract_address = '' OR contract_address IS NULL) AND (token_id = '' OR token_id IS NULL) AND (client_token_id = '' OR client_token_id IS NULL)) OR
(type = 'ERC20' AND contract_address LIKE '0x%' AND (token_id = '' OR token_id IS NULL) AND (client_token_id = '' OR client_token_id IS NULL)) OR
(type = 'ERC721' AND contract_address LIKE '0x%' AND client_token_id IS NOT NULL));
create EXTENSION if not exists citext;

alter table accounts alter column ether_key type citext;

create table if not exists royalties
(
    asset_id            varchar references tokens (id),
    originator_address  varchar(42) not null,
    fee_percentage      integer check (fee_percentage >= 0),
    primary key (asset_id, originator_address)
);

CREATE FUNCTION exists_in_accounts() RETURNS trigger AS $exists_in_accounts$
    BEGIN
        -- If the account doesn't exist we should raise an exception, otherwise proceed.
        IF NOT EXISTS (select 1 from accounts where ether_key = NEW.originator_address) THEN
            RAISE EXCEPTION 'Account % does not exist', NEW.originator_address;
        END IF;
        RETURN NEW;
    END;
$exists_in_accounts$ LANGUAGE plpgsql;

--ensures corresponding accounts exist through trigger on royalties
CREATE TRIGGER account_must_exist BEFORE INSERT OR UPDATE ON royalties
    FOR EACH ROW EXECUTE PROCEDURE exists_in_accounts();

CREATE TYPE fee_type AS ENUM ('royalty', 'ecosystem', 'protocol');

CREATE TYPE stage_type AS ENUM('source', 'proxy', 'dest');

CREATE TABLE IF NOT EXISTS fees
(
    id       bigint REFERENCES orders(id),
    fee_type fee_type,
    amount   NUMERIC(64,0) CHECK (amount >= 0),
    vault_id INTEGER REFERENCES vaults (id),
    stage    stage_type,
    primary key (id, fee_type, vault_id, stage)
);

ALTER TABLE IF EXISTS transactions
    ADD COLUMN log_index numeric(64, 0);

CREATE INDEX transactions_tx_hash_log_index_index ON transactions(tx_hash, log_index);
ALTER TABLE IF EXISTS tokens ADD COLUMN IF NOT EXISTS duplicate BOOLEAN DEFAULT NULL;

CREATE UNIQUE INDEX tokens_erc721_contract_address_token_id_idx ON tokens (contract_address, token_id) WHERE (type = 'ERC721' AND token_id != '' AND duplicate IS NULL);
CREATE UNIQUE INDEX tokens_erc721_contract_address_client_token_id_idx ON tokens (contract_address, client_token_id) WHERE (type = 'ERC721' AND client_token_id != '' AND duplicate IS NULL);
CREATE INDEX tokens_contract_address_type_idx ON tokens (contract_address, type);
CREATE INDEX contracts_contract_address_idx ON contracts (contract_address);

CREATE INDEX IF NOT EXISTS orders_stark_key_index ON orders(stark_key);

CREATE INDEX IF NOT EXISTS orders_buy_vault_id_index ON orders(buy_vault_id);

CREATE INDEX IF NOT EXISTS orders_sell_vault_id_index ON orders(sell_vault_id);

CREATE INDEX IF NOT EXISTS vaults_asset_id_index ON vaults(asset_id);

CREATE TYPE client_enum AS ENUM('immutable', 'rewardpool', 'starkware');

CREATE TABLE IF NOT EXISTS clients
(
    name client_enum primary key,
    address varchar(42)
);

ALTER TABLE transitions ADD fee_paid boolean default false;

create type claim_status as enum (
    'unclaimed', -- User is eligible for claim, but has not yet claimed it
    'claimed', -- User has submitted a claim for their reward
    'paid' -- Claimed amount has been paid out to User
);

/*
Article explaining campaign types
https://immutablex.medium.com/huge-imx-airdrop-for-early-backers-of-immutable-x-the-first-l2-for-nfts-48e731cd4d
 */
create type campaign_type as enum (
    'retrospective', -- Retrospective campaign rewards
    'alpha' -- Alpha campaign rewards
);

create table if not exists claims (
    stark_key varchar(66),
    type campaign_type,
    amount   NUMERIC(64,0) CHECK (amount >= 0),
    token_address varchar(42),
    points varchar,
    expiration_time timestamp,
    status claim_status,
    primary key (stark_key, type)
);

CREATE INDEX IF NOT EXISTS transactions_nonce_index ON transactions(nonce);

-- create trigger to check for unique (stark_key, nonce)
CREATE FUNCTION unique_transactions() RETURNS trigger AS $unique_transactions$
    DECLARE
        target_stark_key TEXT;
    BEGIN
        SELECT stark_key INTO target_stark_key FROM vaults v WHERE v.id = NEW.from_vault_id;
        -- If a transaction (transfer only) with the same (stark_key, nonce) pair
        -- has been created, we should raise an exception, otherwise proceed.
        IF NEW.transaction_type = 'transfer' AND EXISTS (
            SELECT 1 FROM
            vaults v WHERE
            v.id IN (SELECT DISTINCT from_vault_id FROM transactions t  WHERE t.nonce = NEW.nonce)
            AND v.stark_key = target_stark_key
        ) THEN
            RAISE EXCEPTION 'The same transaction with your stark key % and nonce % already exists', target_stark_key, NEW.nonce;
        END IF;
        RETURN NEW;
    END;
$unique_transactions$ LANGUAGE plpgsql;

CREATE TRIGGER unique_transactions BEFORE INSERT OR UPDATE ON transactions
    FOR EACH ROW EXECUTE PROCEDURE unique_transactions();

CREATE INDEX transitions_type_status_idx ON transitions ("type", "status");
CREATE INDEX transactions_from_vault_id_idx ON transactions (from_vault_id);
CREATE INDEX transactions_transaction_type_status_idx ON transactions (transaction_type, "status");

-- Backfill token incentive balance history into the existing balance-change list.
--
-- Token incentive claims already credited user balance before this migration, so
-- this only creates an audit/listing row in redeem_codes and does not touch users.
INSERT INTO redeem_codes (code, type, value, status, used_by, used_at, notes, created_at)
SELECT CONCAT('TI-CLAIM-', c.id),
       'token_incentive',
       c.reward_amount,
       'used',
       c.user_id,
       c.claimed_at,
       CONCAT('Token incentive reward: week ', c.week_start::date, ' ~ ', c.week_end::date, ', tokens=', c.tokens),
       c.claimed_at
FROM token_incentive_claims c
WHERE c.status = 'claimed'
  AND NOT EXISTS (
      SELECT 1
      FROM redeem_codes r
      WHERE r.code = CONCAT('TI-CLAIM-', c.id)
  )
ON CONFLICT (code) DO NOTHING;

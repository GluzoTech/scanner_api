-- Scans captured by the web scanner (scanner.gluzone.dev), where an operator
-- sets a reference model on their phone and then scans barcodes against it.
--
-- Kept separate from `scans`: that table is fed by the hardware agent and is
-- keyed on the device that captured a scan (device_id, vendor_id, product_id,
-- agent_host). A web scan has none of those, and carries a reference_model
-- instead, so folding the two together would leave most columns NULL on every
-- row and make "which source produced this?" a guess.
CREATE TABLE IF NOT EXISTS scan_events (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    reference_model VARCHAR(255) NOT NULL,
    barcode_input TEXT NOT NULL,
    barcode_format VARCHAR(32),
    -- When the phone decoded the barcode, and when the API stored it. Both are
    -- kept because they diverge: a scan queued offline can arrive hours late,
    -- and the phone's clock is the one that reflects when work happened.
    --
    -- TIMESTAMPTZ rather than the bare DATE / TIME pair `scans` uses. The web
    -- scanner sends RFC3339 UTC instants, and an offline queue replaying a
    -- backlog needs an unambiguous ordering that survives a device in another
    -- timezone. Rendering these in IST is a read-time concern.
    scanned_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- The scan log queries by reference model, newest first.
CREATE INDEX IF NOT EXISTS scan_events_reference_model_scanned_at_idx
    ON scan_events (reference_model, scanned_at DESC);

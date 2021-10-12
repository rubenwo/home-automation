CREATE TABLE ids_tradfriids
(
    id         UUID NOT NULL,
    tradfri_id TEXT NOT NULL,

    PRIMARY KEY (id, tradfri_id)
)
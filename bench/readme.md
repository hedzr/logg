# bench and states

`logg/slog` has these states:

- disabled scenes:  <= 1ns | std-slog: 40ns
- withoud fields: ~= 95ns | std-slog: 143ns
- with fields/add fields: ~= 140337ns | std-slog: 148ns

The reason is at serializing attributes.

The tuning wasn't scheduled yet.

```bash
```

TODO.

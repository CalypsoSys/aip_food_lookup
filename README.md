# aip_food_lookup

Go API backend and Flutter migration workspace for AIP Food Lookup.

## Backend

The Go backend lives in `cmd/aip_food_lookup` and serves:

- `GET /search?key=<text>&type=<searchbytextandsound|searchbytext|searchbysound>`
- `POST /suggest`
- `POST /feedback`
- `GET /categories`
- `GET /subcategory?cat=<Allowed|Not Allowed>&sub=<subcategory>`

Food data is stored in `data/allowed` and `data/not_allowed`. Runtime suggestion and feedback files are ignored by git. Feedback currently writes to `data/feedback.jsonl` as plumbing for a future Slack notifier; no Slack webhook or secret is stored in this repository.

Run locally:

```powershell
cd cmd\aip_food_lookup
$env:AIP_DATA_FOLDER='..\..\data'
go run .
```

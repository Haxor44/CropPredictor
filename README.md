# CropPredictor

CropPredictor is a machine learning project that recommends the most suitable
crop to grow based on soil and weather parameters. It combines three pieces:

1. A Jupyter notebook that trains the model from a public crop-recommendation
   dataset.
2. A Django REST API that loads the trained model and returns predictions.
3. A small Go web server that provides a user-friendly HTML form and forwards
   the submitted values to the Django API.

## Model

The model is trained in [`cropoPrediction.ipynb`](cropoPrediction.ipynb) and
predicts a crop label given the following seven soil/weather measurements:

- `N` – nitrogen content in the soil
- `P` – phosphorus content in the soil
- `K` – potassium content in the soil
- `temperature` – ambient temperature (°C)
- `humidity` – relative humidity (%)
- `ph` – soil pH
- `rainfall` – rainfall (mm)

The trained classifier is serialised with `joblib` to
`cropPredictor/cropPredictor.sav` and is loaded at request time by the
Django API.

## Repository layout

- `cropoPrediction.ipynb` – Notebook used to explore the data and train the
  classifier.
- `cropPredictor/` – Django project that exposes the model as a REST API.
  - `cropPredictor/` – Django project settings (`settings.py`, `urls.py`,
    `wsgi.py`, `asgi.py`).
  - `api/` – Django app containing the API views.
    - `views.py` – Defines the `getData` (GET `/`) and `getCrop`
      (POST `/add/`) endpoints.
    - `urls.py` – URL routing for the `api` app.
  - `cropPredictor.sav` / `cropPredictions.sav` – Pickled scikit-learn models
    loaded by `views.getCrop`.
  - `manage.py` – Standard Django management entry-point.
- `crop/` – Go HTTP server that serves the form and calls the Django API.
  - `main.go` – HTTP handler on port `9000`.
  - `predicts.html` – HTML form collecting the seven measurements.

## API

### `GET /`
Smoke-test endpoint. Returns a hard-coded JSON object, for example:

```json
{"name": "Toto", "age": 20}
```

### `POST /add/`
Accepts a JSON body containing the seven measurements in the order
`[P, N, K, temperature, humidity, pH, rainfall]` and returns the crop
predicted by the model.

Example request body:

```json
[
  {
    "measurements": [90, 42, 43, 20, 82, 6, 202]
  }
]
```

Example response:

```json
["rice"]
```

## Running locally

### 1. Django API (port 8000)

Requires Python 3.9+.

```bash
cd cropPredictor
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\\Scripts\\activate
pip install django djangorestframework joblib numpy scikit-learn
python manage.py migrate
python manage.py runserver
```

The API will be available at `http://localhost:8000/`:

- `GET  http://localhost:8000/`       → smoke-test endpoint
- `POST http://localhost:8000/add/`   → crop prediction

### 2. Go web front-end (port 9000)

Requires Go 1.18+.

```bash
cd crop
go run main.go
```

Then open `http://localhost:9000/` in your browser, fill in the form and
submit. The Go server forwards the values to the Django API running on
port 8000 and renders the predicted crop back to the page.

> The Django API URL is currently hard-coded as `http://localhost:8000/add/`
> in `crop/main.go`. Edit the `url` constant in `sendData()` if you run the
> API on a different host or port.

## Retraining the model

Open `cropoPrediction.ipynb` in Jupyter / VS Code, run the cells end-to-end
and save the resulting estimator to `cropPredictor/cropPredictor.sav` with
`joblib.dump`. The Django API will pick up the new file the next time
`views.getCrop` is invoked.

## Notes

- The Django `SECRET_KEY` in `cropPredictor/settings.py` is a development
  key that should be replaced with a secret loaded from the environment
  before deploying anywhere beyond `localhost`.
- `DEBUG = True` is enabled by default and should be disabled in production.

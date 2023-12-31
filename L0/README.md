# wb-course


# Deploy with docker compose

``` shell
$ docker compose --env-file subscriber/.env up
```

Stop the containers with
``` shell
$ docker compose down
# To delete all data run:
$ docker compose down -v
```

Publish any messages __data5.json__ with
```
$ cd pub
$ go run main.go nameOfSubject < data5.json
```

Run GUI app to get order info:
```
$ cd gui
$ python -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
$ streamlit run main.py
```

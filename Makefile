build:
	docker buildx build --platform linux/amd64 -t lucky-dr3w .
	docker tag lucky-dr3w registry.heroku.com/lucky-dr3w/web
deploy:
	docker push registry.heroku.com/lucky-dr3w/web && heroku container:release web --app lucky-dr3w

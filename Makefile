.PHONY: dev-backend dev-frontend clean

dev-backend:
	air

dev-frontend:
	npm --prefix ui run dev

clean:
	rm -rf air/tmp ui/dist
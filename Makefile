.PHONY: docs

test:
	nosetests --with-coverage --cover-package=transiter --rednose -v tests

reset-db:
	python -m transiter.rebuilddb

reset-docs:
	cd docs; rm -r source; rm -r _build; sphinx-apidoc -o source ../transiter; make html

docs:
	cd docs; rm -r _build; make html

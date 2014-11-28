FROM scratch

ADD bin/navyapi /navyapi

CMD ["/navyapi"]

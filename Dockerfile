FROM scratch

COPY imis /imis

ENTRYPOINT [ "./imis" ]
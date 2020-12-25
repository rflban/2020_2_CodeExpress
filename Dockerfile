FROM golang:latest

COPY . /project

WORKDIR /project

RUN make build
RUN mkdir avatars
RUN mkdir album_posters
RUN mkdir artist_avatars
RUN mkdir artist_posters
RUN mkdir playlist_posters
RUN mkdir track_audio

EXPOSE 8085 8086 8087

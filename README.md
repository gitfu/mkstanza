# mkstanza
Generate variant stanza for master.m3u8 file 

## Install
Requires Go and ffprobe
git clone https://github.com/gitfu/mkstanza.git
cd mkstanza
go build mkstanza.go


./mkstanza -i mp3.m3u8

```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
mp3.m3u8
```

./mkstanza -i mp3.m3u8 -u http://example.com

```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
http://example.com/mp3.m3u8
```
./mkstanza  -i audio_and_video.m3u8

```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1994969,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.5"
audio_and_video.m3u8
```

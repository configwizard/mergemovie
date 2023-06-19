# MergeMovie

## About

This is a little app to make downloading videos from websites a breeze. In the simplest for, a video will just be a link to the mp4/avi etc on the site. In that case you can just go to that URL and download that video.

However longer videos or a lot of streaming sites will use (m3u8)[https://en.wikipedia.org/wiki/M3U] format and (transport stream)[https://en.wikipedia.org/wiki/MPEG_transport_stream] files. 

In short an m3u8 file holds the link/path to the ts files. The ts files together make up the whole video.

If you can find these and stitch them together, then you can download the whole video.

## Usage

You can watch the video (here)[https://youtu.be/4s7Q5m4Ma_4] however in short, you just need to find the m3u8 file from the network tab in Chrome/browser (right click, inspect -> network tab), filter by `m3u8` and refresh the page. Look for the m3u8 file that contains a load of .ts files. Copy the m3u8 path and paste into MergeMovie.

Click Download

Choose a file name and location to save it

Done.

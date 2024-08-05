# insert_songs.py
import os
import django

# Set up the Django environment
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'backend.settings')
django.setup()

from music.models import Song
import requests
from django.core.files.base import ContentFile, File
from django.core.files.temp import NamedTemporaryFile

# List of songs to be inserted
songs = [
    {
        "name": "Irreplaceable",
        "author": "Sapio",
        "file": "https://files.freemusicarchive.org/storage-freemusicarchive-org/tracks/BYl1B50WKQCKbiLtECjNYumY18htQTlaqJ5MpUTt.mp3?download=1&name=Sapio%20-%20Irreplaceable.mp3",
        "thumbnail": "https://freemusicarchive.org/image/?file=track_image%2Fvu1EU9H0Pxv67II9pAOvQ28fHdAZppapGOQgYbxY.png&width=290&height=290&type=track",
        "duration": 300
    },
    {
        "name": "Soulmates (Dear Future Wife / Husband)",
        "author": "Tadz",
        "file": "https://files.freemusicarchive.org/storage-freemusicarchive-org/tracks/vGwJswHSC5hO8wmqNvVKpoE9BaehzkkiBWLMASR7.mp3?download=1&name=Tadz%20-%20Soulmates%20%28Dear%20Future%20Wife%20%2F%20Husband%29.mp3",
        "thumbnail": "https://freemusicarchive.org/image/?file=track_image%2FVn4FvKcfzVnLfiQ4zRqarP18oSfQ19o3gG1FfynR.png&width=290&height=290&type=track",
        "duration": 300
    },
    {
        "name": "High School Crush",
        "author": "Tadz",
        "file": "https://files.freemusicarchive.org/storage-freemusicarchive-org/tracks/3Y0b1YV4ePQ0ZwEfMUmDXOtvArQz4Nsoe2rO777W.mp3?download=1&name=Tadz%20-%20High%20School%20Crush.mp3",
        "thumbnail": "https://freemusicarchive.org/image/?file=track_image%2F0dyea9xG9zypHMyiVJ58CyQQWT8Ena2JWZKB0QH0.jpg&width=290&height=290&type=track",
        "duration": 317
    }
]


def download_file(url):
    response = requests.get(url)
    if response.status_code == 200:
        return ContentFile(response.content)
    else:
        return None


for song_data in songs:
    song = Song(
        name=song_data['name'],
        artist=song_data['author'],
        duration=song_data['duration']
    )

    # Download and save the thumbnail
    thumbnail = download_file(song_data['thumbnail'])
    if thumbnail:
        temp_thumb = NamedTemporaryFile(delete=True)
        temp_thumb.write(thumbnail.read())
        temp_thumb.flush()
        song.thumbnail.save(f"{song_data['name']}_thumbnail.jpg", File(temp_thumb))

    # Download and save the song file
    song_file = download_file(song_data['file'])
    if song_file:
        temp_file = NamedTemporaryFile(delete=True)
        temp_file.write(song_file.read())
        temp_file.flush()
        song.file.save(f"{song_data['name']}.mp3", File(temp_file))

    song.save()
    print(f"Inserted {song.name} by {song.artist}")

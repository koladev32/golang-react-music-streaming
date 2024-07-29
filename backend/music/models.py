from django.db import models


class Song(models.Model):
    name = models.CharField(max_length=100)
    artist = models.CharField(max_length=100)
    duration = models.IntegerField()
    thumbnail = models.ImageField(upload_to='images/')
    file = models.FileField(upload_to='music/')
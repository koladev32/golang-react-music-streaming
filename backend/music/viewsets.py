from django.core.cache import cache
from rest_framework import viewsets, permissions, status
from django.utils.decorators import method_decorator
from django.views.decorators.cache import cache_page
from rest_framework.response import Response

from music.models import Song
from music.serializers import SongSerializer


class SongViewSet(viewsets.ModelViewSet):
    queryset = Song.objects.all()
    serializer_class = SongSerializer
    permission_classes = [permissions.AllowAny]

    cache_key = 'song_list_cache'

    @method_decorator(cache_page(60 * 5))  # 5 minutes
    def list(self, request, *args, **kwargs):
        return super(SongViewSet, self).list(request, *args, **kwargs)

    def create(self, request, *args, **kwargs):
        response = super(SongViewSet, self).create(request, *args, **kwargs)
        if response.status_code == status.HTTP_201_CREATED:
            # Invalidate cache on create
            cache.delete(self.cache_key)
        return response

    def update(self, request, *args, **kwargs):
        response = super(SongViewSet, self).update(request, *args, **kwargs)
        if response.status_code == status.HTTP_200_OK:
            # Invalidate cache on update
            cache.delete(self.cache_key)
        return response

    def destroy(self, request, *args, **kwargs):
        response = super(SongViewSet, self).destroy(request, *args, **kwargs)
        if response.status_code == status.HTTP_204_NO_CONTENT:
            # Invalidate cache on delete
            cache.delete(self.cache_key)
        return response

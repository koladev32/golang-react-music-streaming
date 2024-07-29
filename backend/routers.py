from rest_framework.routers import SimpleRouter

from music.viewsets import SongViewSet

router = SimpleRouter()

router.register('songs', SongViewSet)

urlpatterns = router.urls
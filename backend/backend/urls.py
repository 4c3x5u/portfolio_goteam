from django.contrib import admin
from django.urls import path
from main.api_views.api_auth import register, login
from main.api_views.api_boards import boards
from main.api_views.api_tasks import tasks

urlpatterns = [
    path('register/', register, name='register'),
    path('login/', login, name='login'),
    path('boards/', boards, name='board'),
    path('tasks/', tasks, name='tasks'),

    path('admin/', admin.site.urls)
]

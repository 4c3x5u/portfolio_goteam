from django.contrib import admin
from django.urls import path
from main.api.api_auth import register, login, verify_token
from main.api.api_boards import boards
from main.api.api_tasks import tasks
from main.api.api_subtasks import subtasks

urlpatterns = [
    path('register/', register, name='register'),
    path('login/', login, name='login'),
    path('verify-token/', verify_token, name='verifytoken'),
    path('boards/', boards, name='board'),
    path('tasks/', tasks, name='tasks'),
    path('subtasks/', subtasks, name='subtasks'),

    path('admin/', admin.site.urls)
]

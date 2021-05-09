from django.urls import path
from main.api.registerapi import Register
from main.api.loginapi import Login
from main.api.api_users import users
from main.api.api_boards import Boards
from main.api.api_columns import columns
from main.api.api_tasks import tasks
from main.api.subtasksapi import Subtasks
from main.api.api_clientstate import ClientState

urlpatterns = [
    path('register/', Register.as_view(), name='register'),
    path('login/', Login.as_view(), name='login'),
    path('users/', users, name='users'),
    # path('boards/', boards, name='boards'),
    path('boards/', Boards.as_view(), name='boards'),
    path('columns/', columns, name='columns'),
    path('tasks/', tasks, name='tasks'),
    path('subtasks/', Subtasks.as_view(), name='subtasks'),
    path('client-state/', ClientState.as_view(), name='clientstate')
]

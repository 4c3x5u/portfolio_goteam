from django.urls import path
from main.api.api_auth import Register, login, verify_token
from main.api.api_users import users
from main.api.api_teams import teams
from main.api.api_boards import boards
from main.api.api_columns import columns
from main.api.api_tasks import tasks
from main.api.api_subtasks import Subtasks
from main.api.api_clientstate import client_state

urlpatterns = [
    path('verify-token/', verify_token, name='verifytoken'),
    path('register/', Register.as_view(), name='register'),
    path('login/', login, name='login'),
    path('users/', users, name='users'),
    path('teams/', teams, name='teams'),
    path('boards/', boards, name='boards'),
    path('columns/', columns, name='columns'),
    path('tasks/', tasks, name='tasks'),
    path('subtasks/', Subtasks.as_view(), name='subtasks'),
    path('client-state/', client_state, name='clientstate')
]

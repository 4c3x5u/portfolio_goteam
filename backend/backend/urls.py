from django.urls import path
from main.views.api_register import RegisterApi
from main.views.api_login import LoginApi
from main.views.api_users import UsersApi
from main.views.api_boards import BoardsApi
from main.views.api_columns import ColumnsApi
from main.views.api_tasks import TasksApi
from main.views.api_subtasks import SubtasksApi
from main.views.api_client_state import ClientStateApi

urlpatterns = [
    path('register/', RegisterApi.as_view(), name='register'),
    path('login/', LoginApi.as_view(), name='login'),
    path('users/', UsersApi.as_view(), name='users'),
    path('boards/', BoardsApi.as_view(), name='boards'),
    path('columns/', ColumnsApi.as_view(), name='columns'),
    path('tasks/', TasksApi.as_view(), name='tasks'),
    path('subtasks/', SubtasksApi.as_view(), name='subtasks'),
    path('client-state/', ClientStateApi.as_view(), name='client-state')
]

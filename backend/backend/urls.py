from django.urls import path
from main.views.api_register import RegisterApiView
from main.views.api_login import LoginApiView
from main.views.api_users import UsersApiView
from main.views.api_boards import BoardsApiView
from main.views.api_columns import ColumnsApiView
from main.views.api_tasks import TasksApiView
from main.views.api_subtasks import SubtasksApiView
from main.views.api_client_state import ClientStateApiView

urlpatterns = [
    path('register/', RegisterApiView.as_view(), name='register'),
    path('login/', LoginApiView.as_view(), name='login'),
    path('users/', UsersApiView.as_view(), name='users'),
    path('boards/', BoardsApiView.as_view(), name='boards'),
    path('columns/', ColumnsApiView.as_view(), name='columns'),
    path('tasks/', TasksApiView.as_view(), name='tasks'),
    path('subtasks/', SubtasksApiView.as_view(), name='subtasks'),
    path('client-state/', ClientStateApiView.as_view(), name='client-state')
]

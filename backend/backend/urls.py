from django.urls import path
from main.api_views.register import RegisterAPIView
from main.api_views.login import LoginAPIView
from main.api_views.users import UsersAPIView
from main.api_views.boards import BoardsAPIView
from main.api_views.columns import ColumnsAPIView
from main.api_views.api_tasks import TasksAPIView
from main.api_views.subtasks import SubtasksAPIView
from main.api_views.clientstate import ClientStateAPIView

urlpatterns = [
    path('register/', RegisterAPIView.as_view(), name='register'),
    path('login/', LoginAPIView.as_view(), name='login'),
    path('users/', UsersAPIView.as_view(), name='users'),
    path('boards/', BoardsAPIView.as_view(), name='boards'),
    path('columns/', ColumnsAPIView.as_view(), name='columns'),
    path('tasks/', TasksAPIView.as_view(), name='tasks'),
    path('subtasks/', SubtasksAPIView.as_view(), name='subtasks'),
    path('client-state/', ClientStateAPIView.as_view(), name='client-state')
]

from django.urls import path
from main.api_views.registerapiview import RegisterAPIView
from main.api_views.loginapiview import LoginAPIView
from main.api_views.usersapiview import UsersAPIView
from main.api_views.boardsapiview import BoardsAPIView
from main.api_views.columnsapiview import ColumnsAPIView
from main.api_views.api_tasks import TasksAPIView
from main.api_views.subtasksapiview import SubtasksAPIView
from main.api_views.clientstateapiview import ClientStateAPIView

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

from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status

from ..models import Task
from ..serializers.task.ser_task_update import UpdateTaskSerializer
from ..validation.val_auth import \
    authenticate, authorize, not_authenticated_response
from ..validation.val_task import validate_task_id
from ..serializers.task.ser_task_create import CreateTaskSerializer


class TasksApiView(APIView):
    @staticmethod
    def post(request):
        serializer = CreateTaskSerializer(
            data={'column': request.data.get('column'),
                  'title': request.data.get('title'),
                  'description': request.data.get('description'),
                  'subtasks': request.data.get('subtasks'),
                  'auth_user': request.META.get('HTTP_AUTH_USER'),
                  'auth_token': request.META.get('HTTP_AUTH_TOKEN')}
        )
        if serializer.is_valid():
            serializer.create(serializer.validated_data)
            return Response(serializer.data, status.HTTP_201_CREATED)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def patch(request):
        data = {'task': request.query_params.get('id'),
                'auth_user': request.META.get('HTTP_AUTH_USER'),
                'auth_token': request.META.get('HTTP_AUTH_TOKEN')}
        for key, value in request.data.items():
            data[key] = value
        serializer = UpdateTaskSerializer(data=data)
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)


@api_view(['PATCH', 'DELETE'])
def tasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'DELETE':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        task_id = request.query_params.get('id')

        validation_response = validate_task_id(task_id)
        if validation_response:
            return validation_response

        try:
            task = Task.objects.select_related(
                'column',
                'column__board',
            ).get(id=task_id)
        except Task.DoesNotExist:
            return Response({
                'task_id': ErrorDetail(string='Task not found.',
                                       code='not_found')
            }, 404)

        if task.column.board.team_id != user.team_id:
            return not_authenticated_response

        task.delete()

        return Response({
            'msg': 'Task deleted successfully.',
            'id': task_id,
        })

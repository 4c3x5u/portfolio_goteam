from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status
from ..models import Task, Column
from main.serializers.task.taskserializer import TaskSerializer
from main.serializers.subtask.subtaskserializer import SubtaskSerializer
from ..validation.auth import \
    authenticate, authorize, not_authenticated_response
from ..validation.column import validate_column_id
from ..validation.task import validate_task_id
from ..serializers.task.createtaskserializer import CreateTaskSerializer


class TasksAPIView(APIView):
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


@api_view(['PATCH', 'DELETE'])
def tasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        task_id = request.query_params.get('id')
        validation_response = validate_task_id(task_id)
        if validation_response:
            return validation_response

        task = Task.objects.select_related(
            'column',
            'column__board'
        ).prefetch_related(
            'subtask_set'
        ).get(id=task_id)

        if task.column.board.team_id != user.team_id:
            return not_authenticated_response

        if 'title' in request.data.keys() and not request.data.get('title'):
            return Response({
                'title': ErrorDetail(string='Task title cannot be empty.',
                                     code='blank')
            }, 400)

        order = request.data.get('order')
        if 'order' in request.data.keys() and (order == '' or order is None):
            return Response({
                'order': ErrorDetail(string='Task order cannot be empty.',
                                     code='blank')
            }, 400)

        if 'column' in request.data.keys():
            column_id = request.data.get('column')
            validation_response = validate_column_id(column_id)
            if validation_response:
                return validation_response

        subtasks = request.data.pop(
            'subtasks'
        ) if 'subtasks' in request.data.keys() else None

        # update tasks
        task_serializer = TaskSerializer(task, data=request.data, partial=True)
        if not task_serializer.is_valid():
            return Response(task_serializer.errors, 400)
        task = task_serializer.save()

        # update subtasks
        if subtasks:
            task.subtask_set.all().delete()
            subtasks_data = [
                {
                    'title': subtask['title'],
                    'order': subtask['order'],
                    'task': task.id,
                    'done': subtask['done']
                } for subtask in subtasks
            ]
            subtask_serializer = SubtaskSerializer(data=subtasks_data,
                                                   many=True)
            if not subtask_serializer.is_valid():
                return Response({'subtasks': subtask_serializer.errors}, 400)
            subtask_serializer.save()

        return Response({
            'msg': 'Task update successful.',
            'id': task.id
        }, 200)

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

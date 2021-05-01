from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask, Task
from ..serializers.ser_subtask import SubtaskSerializer
from ..validation.val_auth import \
    authenticate, authorize, not_authenticated_response
from ..validation.val_task import validate_task_id


@api_view(['GET', 'PATCH'])
def subtasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    # not in use – maintained for demonstration purposes
    if request.method == 'GET':
        task_id = request.query_params.get('task_id')

        validation_response = validate_task_id(task_id)
        if validation_response:
            return validation_response

        try:
            task = Task.objects.select_related(
                'column',
                'column__board',
            ).prefetch_related(
                'subtask_set'
            ).get(id=task_id)
        except Task.DoesNotExist:
            return Response({
                'task_id': ErrorDetail(string='Task not found.',
                                       code='not_found')
            }, 404)

        if task.column.board.team_id != user.team_id:
            return not_authenticated_response

        return Response({
            'subtasks': [
                {
                    'id': subtask.id,
                    'order': subtask.order,
                    'title': subtask.title,
                    'done': subtask.done
                } for subtask in task.subtask_set.all()
            ]
        }, 200)

    if request.method == 'PATCH':
        subtask_id = request.query_params.get('id')
        if not subtask_id:
            return Response({
                'id': ErrorDetail(string='Subtask ID cannot be empty.',
                                  code='blank')
            }, 400)

        try:
            subtask = Subtask.objects.select_related(
                'task',
                'task__user',
                'task__column__board'
            ).get(id=subtask_id)
        except Subtask.DoesNotExist:
            return Response({
                'id': ErrorDetail(string='Subtask not found.',
                                  code='not_found')
            })

        authorization_response = authorize(username)

        # if the user is NOT admin and is NOT assigned to this task...
        if authorization_response and subtask.task.user != user:
            return authorization_response

        if subtask.task.column.board.team_id != user.team_id:
            return not_authenticated_response

        if not request.data:
            return Response({
                'data': ErrorDetail(string='Data cannot be empty.',
                                    code='blank')
            }, 400)

        if 'title' in request.data.keys() and not request.data.get('title'):
            return Response({
                'title': ErrorDetail(string='Title cannot be empty.',
                                     code='blank')
            }, 400)

        done = request.data.get('done')
        if 'done' in request.data.keys() and (done == '' or done is None):
            return Response({
                'done': ErrorDetail(string='Done cannot be empty.',
                                    code='blank')
            }, 400)

        order = request.data.get('order')
        if 'order' in request.data.keys() and (order == '' or order is None):
            return Response({
                'order': ErrorDetail(string='Order cannot be empty.',
                                     code='blank')
            }, 400)

        serializer = SubtaskSerializer(subtask,
                                       data=request.data,
                                       partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)

        subtask = serializer.save()
        return Response({
            'msg': 'Subtask update successful.',
            'id': subtask.id
        }, 200)


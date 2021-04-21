from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask, Task
from ..serializers.ser_subtask import SubtaskSerializer
from ..util import authenticate, authorize, not_authenticated_response


@api_view(['GET', 'PATCH'])
def subtasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    team_id, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        task_id = request.query_params.get('task_id')

        if not task_id:
            return Response({
                'task_id': ErrorDetail(string='Task ID cannot be empty.',
                                       code='blank')
            }, 400)

        try:
            int(task_id)
        except ValueError:
            return Response({
                'task_id': ErrorDetail(string='Task ID must be a number.',
                                       code='invalid')
            }, 400)

        try:
            task = Task.objects.get(id=task_id)
        except Task.DoesNotExist:
            return Response({
                'task_id': ErrorDetail(string='Task not found.',
                                       code='not_found')
            }, 404)

        if task.column.board.team_id != team_id:
            return not_authenticated_response

        task_subtasks = Subtask.objects.filter(task_id=task_id)
        serializer = SubtaskSerializer(task_subtasks, many=True)
        return Response({
            'subtasks': list(
                map(
                    lambda st: {
                        'id': st['id'],
                        'order': st['order'],
                        'title': st['title'],
                        'done': st['done']
                    },
                    serializer.data
                )
            )
        }, 200)

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        subtask_id = request.query_params.get('id')
        if not subtask_id:
            return Response({
                'id': ErrorDetail(string='Subtask ID cannot be empty.',
                                  code='blank')
            }, 400)

        try:
            subtask = Subtask.objects.get(id=subtask_id)
        except Subtask.DoesNotExist:
            return Response ({
                'id': ErrorDetail(string='Subtask not found.',
                                  code='not_found')
            })

        if subtask.task.column.board.team.id != team_id:
            return not_authenticated_response

        data = request.data

        if not data:
            return Response({
                'data': ErrorDetail(string='Data cannot be empty.',
                                    code='blank')
            }, 400)

        if 'title' in list(data.keys()) and not data.get('title'):
            return Response({
                'title': ErrorDetail(string='Title cannot be empty.',
                                          code='blank')
            }, 400)

        done = data.get('done')
        if 'done' in list(data.keys()) and (done == '' or done is None):
            return Response({
                'done': ErrorDetail(string='Done cannot be empty.',
                                         code='blank')
            }, 400)

        order = data.get('order')
        if 'order' in list(data.keys()) and (order == '' or order is None):
            return Response({
                'order': ErrorDetail(string='Order cannot be empty.',
                                          code='blank')
            }, 400)

        serializer = SubtaskSerializer(Subtask.objects.get(id=subtask_id),
                                       data=data,
                                       partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)

        subtask = serializer.save()
        return Response({
            'msg': 'Subtask update successful.',
            'id': subtask.id
        }, 200)


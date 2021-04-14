from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Task
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer
from ..util import authenticate, authorize


@api_view(['POST', 'PATCH', 'GET'])
def tasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        column_id = request.query_params.get('column_id')

        if not column_id:
            return Response({
                'column_id': ErrorDetail(string='Column ID cannot be empty.',
                                         code='blank')
            }, 400)

        try:
            int(column_id)
        except ValueError:
            return Response({
                'column_id': ErrorDetail(string='Column ID must be a number.',
                                         code='invalid')
            }, 400)

        try:
            Column.objects.get(id=column_id)
        except Column.DoesNotExist:
            return Response({
                'column_id': ErrorDetail(string='Column not found.',
                                         code='not_found')
            }, 404)

        column_tasks = Task.objects.filter(column_id=column_id)
        serializer = TaskSerializer(column_tasks, many=True)
        return Response({
            'tasks': list(map(
                lambda t: {'id': t['id'],
                           'order': t['order'],
                           'title': t['title'],
                           'description': t['description']}
                , serializer.data
            ))
        }, 200)

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        column_id = request.data.get('column')
        if not column_id:
            return Response({
                'column': ErrorDetail(string='Column cannot be empty.',
                                      code='blank')
            }, 400)

        task_serializer = TaskSerializer(
            data={'title': request.data.get('title'),
                  'description': request.data.get('description'),
                  'order': Column.objects.filter(id=column_id).count() + 1,
                  'column': request.data.get('column')}
        )
        if not task_serializer.is_valid():
            return Response(task_serializer.errors, 400)
        task = task_serializer.save()

        subtasks = request.data.get('subtasks')
        if subtasks:
            for i, subtask in enumerate(subtasks):
                subtask_serializer = SubtaskSerializer(
                    data={'title': subtask.get('title'),
                          'order': i,
                          'task': task.id}
                )
                if not subtask_serializer.is_valid():
                    task.delete()
                    return Response({
                        'subtask': subtask_serializer.errors
                    }, 400)
                subtask_serializer.save()

        return Response({
            'msg': 'Task creation successful.',
            'task_id': task.id
        }, 201)

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        task_id = request.data.get('id')
        data = request.data.get('data')

        if 'title' in list(data.keys()) and not data.get('title'):
            return Response({
                'data.title': ErrorDetail(string='Task title cannot be empty.',
                                          code='blank')
            }, 400)

        order = data.get('order')
        if 'order' in list(data.keys()) and (order == '' or order is None):
            return Response({
                'data.order': ErrorDetail(string='Task order cannot be empty.',
                                          code='blank')
            }, 400)

        if 'column' in list(data.keys()):
            column_id = data.get('column')
            if not column_id:
                return Response({
                    'data.column': ErrorDetail(string='Task column cannot be '
                                                      'empty.',
                                               code='blank')
                }, 400)

            try:
                Column.objects.get(id=column_id)
            except Column.DoesNotExist:
                return Response({
                    'data.column': ErrorDetail(string='Invalid column id.',
                                               code='invalid')
                }, 400)

        serializer = TaskSerializer(Task.objects.get(id=task_id),
                                    data=data,
                                    partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)

        task = serializer.save()
        return Response({
            'msg': 'Task update successful.',
            'id': task.id
        }, 200)

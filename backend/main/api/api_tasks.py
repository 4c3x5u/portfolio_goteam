from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Task, User
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer
from .util import authenticate


@api_view(['POST', 'PATCH'])
def tasks(request):
    # validate username
    auth_response = authenticate(request)
    if auth_response:
        return auth_response

    if request.method == 'POST':
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
                    return Response({'subtask': subtask_serializer.errors}, 400)
                subtask_serializer.save()

        return Response({
            'msg': 'Task creation successful.',
            'task_id': task.id
        }, 201)

    if request.method == 'PATCH':
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



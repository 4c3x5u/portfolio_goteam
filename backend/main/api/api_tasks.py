from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Task, Subtask
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer
from ..util import (
    authenticate, authorize, not_authenticated_response, validate_column_id,
    validate_task_id
)


@api_view(['GET', 'POST', 'PATCH', 'DELETE'])
def tasks(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    team_id, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        column_id = request.query_params.get('column_id')

        column, validation_response = validate_column_id(column_id)
        if validation_response:
            return validation_response

        if column.board.team_id != team_id:
            return not_authenticated_response

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
        column, validation_response = validate_column_id(column_id)
        if validation_response:
            return validation_response

        if column.board.team.id != team_id:
            return not_authenticated_response

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
                    data={'title': subtask,
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

        task_id = request.query_params.get('id')

        task, validation_response = validate_task_id(task_id)
        if validation_response:
            return validation_response

        if task.column.board.team.id != team_id:
            return not_authenticated_response

        data = request.data

        if 'title' in list(data.keys()) and not data.get('title'):
            return Response({
                'title': ErrorDetail(string='Task title cannot be empty.',
                                     code='blank')
            }, 400)

        order = data.get('order')
        if 'order' in list(data.keys()) and (order == '' or order is None):
            return Response({
                'order': ErrorDetail(string='Task order cannot be empty.',
                                     code='blank')
            }, 400)

        if 'column' in list(data.keys()):
            column_id = data.get('column')
            _, validation_response = validate_column_id(column_id)
            if validation_response:
                return validation_response

        subtasks = data.get('subtasks')
        'subtask' in list(data.keys()) and data.pop('subtasks')

        task_serializer = TaskSerializer(Task.objects.get(id=task_id),
                                         data=data,
                                         partial=True)
        if not task_serializer.is_valid():
            return Response(task_serializer.errors, 400)
        task = task_serializer.save()

        Subtask.objects.filter(task_id=task.id).delete()

        if subtasks:
            for i, subtask in enumerate(subtasks):
                subtask_serializer = SubtaskSerializer(
                    data={'title': subtask['title'],
                          'order': i,
                          'task': task.id,
                          'done': subtask['done']}
                )
                if not subtask_serializer.is_valid():
                    return Response({
                        'subtasks': subtask_serializer.errors
                    }, 400)
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

        task, validation_response = validate_task_id(task_id)
        if validation_response:
            return validation_response

        if task.column.board.team.id != team_id:
            return not_authenticated_response

        task.delete()

        return Response({
            'msg': 'Task deleted successfully.',
            'id': task_id,
        })

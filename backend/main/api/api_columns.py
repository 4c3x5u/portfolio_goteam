from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Task
from ..serializers.ser_task import TaskSerializer
from ..validation.val_auth import authenticate, authorize, \
    not_authenticated_response
from ..validation.val_column import validate_column_id


@api_view(['GET', 'PATCH'])
def columns(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'PATCH':
        authorization_response = authorize(username)

        column_id = request.query_params.get('id')
        validation_response = validate_column_id(column_id)
        if validation_response:
            return validation_response

        try:
            column = Column.objects.select_related('board').get(id=column_id)
        except Column.DoesNotExist:
            return Response({
                'column_id': ErrorDetail(string='Column not found.',
                                         code='not_found')
            }, 404)

        if column.board.team_id != user.team_id:
            return not_authenticated_response

        # retrieve tasks to avoid a call to the DB for each task
        tasks = Task.objects.filter(column__board_id=column.board_id)

        for task in request.data:
            try:
                task_id = task.pop('id')
            except KeyError:
                return Response({
                    'task.id': ErrorDetail(string='Task ID cannot be empty.',
                                           code='blank')
                }, 400)

            existing_task = tasks.get(id=task_id)

            if authorization_response \
                    and task['user'] != user.username \
                    and column.id != existing_task.column_id:
                return authorization_response

            serializer = TaskSerializer(existing_task,
                                        data={**task, 'column': column.id},
                                        partial=True)
            if not serializer.is_valid():
                return Response(serializer.errors, 400)

            serializer.save()

        return Response({
            'msg': 'Column and all its tasks updated successfully.',
            'id': column.id,
        }, 200)

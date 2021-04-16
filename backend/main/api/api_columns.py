from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..util import authenticate, authorize
from ..models import Column, Board, Task
from ..serializers.ser_column import ColumnSerializer
from ..serializers.ser_task import TaskSerializer


@api_view(['GET', 'PATCH'])
def columns(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        board_id = request.query_params.get('board_id')

        if not board_id:
            return Response({
                'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                        code='blank')
            }, 400)

        try:
            int(board_id)
        except ValueError:
            return Response({
                'board_id': ErrorDetail(string='Board ID must be a number.',
                                        code='invalid')
            }, 400)

        try:
            Board.objects.get(id=board_id)
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        board_columns = Column.objects.filter(board_id=board_id)

        if not board_columns:
            board_columns = [
                Column.objects.create(
                    order=i,
                    board_id=board_id
                ) for i in range(0, 4)
            ]

        serializer = ColumnSerializer(board_columns, many=True)

        return Response({
            'columns': list(
                map(lambda column: {'id': column['id'], 'order': column['order']},
                    serializer.data)
            )
        }, 200)

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        column_id = request.query_params.get('id')

        if not column_id:
            return Response({
                'id': ErrorDetail(string='Column ID cannot be empty.',
                                  code='blank')
            }, 400)

        column = Column.objects.get(id=column_id)

        tasks = request.data

        for task in tasks:
            try:
                task_id = task.pop('id')
            except KeyError:
                return Response({
                    'task.id': ErrorDetail(string='Task ID cannot be empty.',
                                           code='blank')
                }, 400)

            serializer = TaskSerializer(Task.objects.get(id=task_id),
                                        data={**task, 'column': column.id},
                                        partial=True)
            if not serializer.is_valid():
                return Response(serializer.errors, 400)

            serializer.save()

        return Response({
            'msg': 'Column and all its tasks updated successfully.',
            'id': column.id,
        }, 200)

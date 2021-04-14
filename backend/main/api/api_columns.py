from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..util import authenticate
from ..models import Column, Board
from ..serializers.ser_column import ColumnSerializer


@api_view(['GET'])
def columns(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

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
        Board.objects.filter(id=board_id)
    except Board.DoesNotExist:
        return Response({
            'board_id': ErrorDetail(string='Board not found.',
                                    code='not_found')
        }, 404)

    board_columns = Column.objects.filter(board_id=board_id)
    serializer = ColumnSerializer(board_columns, many=True)

    return Response({
        'columns': list(
            map(lambda column: {'id': column['id'], 'order': column['order']},
                serializer.data)
        )
    }, 200)

from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status

from ..models import Task
from ..serializers.task.ser_task_create import CreateTaskSerializer
from ..serializers.task.ser_task_update import UpdateTaskSerializer
from ..serializers.task.ser_task_delete import DeleteTaskSerializer
from ..validation.val_auth import \
    authenticate, authorize, not_authenticated_response
from ..validation.val_task import validate_task_id


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

    @staticmethod
    def delete(request):
        serializer = DeleteTaskSerializer(data={
            'task': request.query_params.get('id'),
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.delete()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

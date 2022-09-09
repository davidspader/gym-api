<?php

namespace App\Http\Controllers;

use App\Http\Requests\API\ExercisesRequest;
use App\Models\Exercise;
use App\Models\User;
use Illuminate\Http\Request;

class ExerciseController extends Controller
{
    private $user;

    public function __construct(private Exercise $exercise)
    {
        $this->user = auth('sanctum')->user();

    }

    public function index()
    {
        return response()->json($this->exercise->all()->where('user_id', $this->user['id']));
    }

    public function show($id)
    {
        $exercise = $this->exercise->find($id);

        $this->authorize('checksUserId', $exercise);

        return response()->json($exercise);
    }

    public function store(ExercisesRequest $request)
    {
        $exercise = $this->exercise->create($request->all());

        $this->authorize('checksUserId', $exercise);

        return response()->json($exercise, 201);
    }
}

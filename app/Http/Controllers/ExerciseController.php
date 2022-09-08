<?php

namespace App\Http\Controllers;

use App\Models\Exercise;
use Illuminate\Http\Request;

class ExerciseController extends Controller
{
    public function __construct(private Exercise $exercise)
    {
    }

    public function index()
    {
        $user = auth('sanctum')->user();
        return response()->json($this->exercise->all()->where('user_id', $user['id']));
    }
}

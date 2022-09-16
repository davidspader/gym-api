<?php

namespace App\Http\Controllers;

use App\Models\Workout;
use Illuminate\Http\Request;

class WorkoutController extends Controller
{
    private $user;

    public function __construct(private Workout $workout)
    {
        $this->user = auth('sanctum')->user();
    }

    public function index()
    {
        return response()->json($this->workout->all()->where('user_id', $this->user['id']));
    }
}

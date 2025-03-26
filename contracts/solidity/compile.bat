@echo off

:: Exit on any error
setlocal enabledelayedexpansion

:: Set default flag values
set no_deps=false
set input_file=

:: Parse arguments
:parse_args
if "%~1" == "--no-deps" (
    set no_deps=true
    shift
    goto parse_args
) else if "%~1" neq "" (
    if "%~x1" == ".sol" (
        set input_file=%~1
    )
    shift
    goto parse_args
)

:: Check if an input file was provided
if "%input_file%" == "" (
    echo Error: Please provide a Solidity file name.
    exit /b 1
)

:: Extract the filename without extension
for %%f in ("%input_file%") do set input_file_name=%%~nxf
set input_file_name=%input_file_name:.sol=%

:: Extract the directory of the input file
for %%f in ("%input_file%") do set input_dir=%%~dpf
set input_dir=%input_dir:~0,-1%

:: Set the output directory
set output_dir=.temp

:: Create the output directory if it doesn't exist
if not exist "%output_dir%" mkdir "%output_dir%"

:: Compile the Solidity file and output to the specified directory
:: Include node_modules in the include path
solcjs "%input_file%" --bin --abi --optimize --base-path . --include-path node_modules -o "%output_dir%"

:: Rename the files
for %%f in ("%output_dir%\*_sol_*.abi" "%output_dir%\*_sol_*.bin") do (
    set filename=%%~nxf
    set filename=!filename:_sol_=!
    set filename=!filename:.abi=!
    set filename=!filename:.bin=!

    :: Determine the extension
    set extension=%%~xf

    :: Rename the file
    ren "%%f" "!filename!%%~xg"
)

:: Move the renamed files to the input directory
if "%no_deps%" == "true" (
    :: Only move files for the target contract
    move /y "%output_dir%\%input_file_name%.abi" "%output_dir%\%input_file_name%.bin" "%input_dir%"
) else (
    :: Move all files
    move /y "%output_dir%\*.abi" "%output_dir%\*.bin" "%input_dir%"
)

:: Delete the temporary output directory
rmdir /s /q "%output_dir%"
